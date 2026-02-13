import type { Context } from "../context";

export type CommandHandler = (
  ctx: Context,
  args: string[]
) => Promise<string[]> | string[];

export class Host {
  private commands = new Map<string, CommandHandler>();

  registerCommand(name: string, handler: CommandHandler, aliases?: string[]) {
    if (this.commands.has(name)) {
      throw new Error(`command already registered: ${name}`);
    }
    this.commands.set(name, handler);
    aliases?.forEach((alias) => this.commands.set(alias, handler))
  }

  async execute(ctx: Context, name: string, args: string[]) {
    const cmd = this.commands.get(name);
    if (!cmd) throw new Error(`unknown command: ${name}`);
    return await cmd(ctx, args);
  }

  getSlice<T>(ctx: Context, pluginId: string, defaultState: T): T {
      if (!ctx.pluginData.has(pluginId)) {
        ctx.pluginData.set(pluginId, defaultState);
      }
      return ctx.pluginData.get(pluginId) as T;
    }

  listCommands() {
    return [...this.commands.keys()];
  }
}
