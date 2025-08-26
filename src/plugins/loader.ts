import { defaultDiagnostics, Diagnostics } from "@/diagnostics";
import { register as regRegister, getSpec, type CommandSpec } from "@/registry";
import { CommandOrigin } from "@/types";
import z from "zod";

export type PluginContext = {
  supportedApiVersion: number;
  projectRoot: string;
  diagnostics: Diagnostics;
};

export type ScopedRegister = {
  command: <TSchema extends z.ZodTypeAny>(
    spec: Omit<CommandSpec<TSchema>, "origin">,
  ) => void;
};

export type PluginModule = {
  name: string;
  version: string;
  apiVersion: number;
  init?: (ctx: PluginContext) => void | Promise<void>;
  register: (
    registerFn: ScopedRegister,
    ctx: PluginContext,
  ) => void | Promise<void>;
  finalize?: (ctx: PluginContext) => void | Promise<void>;
};

export class PluginLoader {
  constructor(
    public supportedApiVersion: number,
    private opts: {
      projectRoot?: string;
      diagnostics?: Diagnostics;
      conflictPolicy?: "error" | "warn";
    } = {},
  ) {}

  private loaded: { plugin: PluginModule; origin: CommandOrigin }[] = [];

  finalize() {}

  public async loadAll() {
    await this.loadCore();
  }

  private async loadCore() {
    const core: PluginModule = await import("@/plugins/core/index.ts");
    const origin: CommandOrigin = {
      name: core.name,
      version: core.version,
      source: "core",
      pluginId: `${core.name}@${core.version}`,
    };

    const ctx = this.makeCtx();

    // TODO: validate manifest of plugin

    try {
      if (core.init) await core.init(ctx);
      const scoped = this.createScopedRegistry(origin, ctx);
      await core.register(scoped, ctx);
      this.loaded.push({ plugin: core, origin });
    } catch (e) {
      ctx.diagnostics.report({
        level: "error",
        code: "plugin/register-error",
        message: String(e instanceof Error ? e.message : e),
        plugin: origin.pluginId,
      });
    }
  }
  private createScopedRegistry(
    origin: CommandOrigin,
    ctx: PluginContext,
  ): ScopedRegister {
    const conflictPolicy = this.opts.conflictPolicy ?? "warn";
    return {
      command: <TSchema extends z.ZodTypeAny>(
        spec: Omit<CommandSpec<TSchema>, "origin">,
      ) => {
        // basic validation; check for name and set kind
        if (!spec?.name || (spec.kind !== "block" && spec.kind !== "inline")) {
          ctx.diagnostics.report({
            level: "error",
            code: "plugin/spec-invalid",
            message: `Invalid command spec for '${spec?.name ?? "<unnamed>"}' caused by name or kind`,
            plugin: origin.pluginId,
          });
          return;
        }
        if (!spec.schema || "safeParse" in spec.schema) {
          ctx.diagnostics.report({
            level: "error",
            code: "plugin/spec-invalid",
            message: `Invalid command spec for '${spec?.name ?? "<unnamed>"} caused by schema'`,
            plugin: origin.pluginId,
          });
          return;
        }
        // validation; check for ir transform function
        if (typeof spec.transform !== "function") {
          ctx.diagnostics.report({
            level: "error",
            code: "plugin/spec-invalid-transform",
            message: `Command '${spec.name}' missing transform()`,
            plugin: origin.pluginId,
          });
        }

        // TODO: implement conflict handling with conflictPolicy
        const fullSpec: CommandSpec = { ...(spec as any), origin };
        regRegister(fullSpec);
      },
    };
  }

  private makeCtx(): PluginContext {
    return {
      supportedApiVersion: this.supportedApiVersion,
      projectRoot: this.opts.projectRoot ?? process.cwd(),
      diagnostics: this.opts.diagnostics ?? defaultDiagnostics(),
    };
  }
}
