import { z } from "zod";
import type { PluginContext, ScopedRegister } from "@/plugins/loader";

export const name = "core";
export const version = "v0.1.0";
export const apiVersion = 1;

export async function init(_ctx: PluginContext) {}

export async function register(
  registerFn: ScopedRegister,
  _ctx: PluginContext,
) {
  registerFn.command({
    name: "heading",
    kind: "block",
    schema: z.tuple([z.coerce.number().int().min(1).max(6), z.string()]),
    transform: (_n, [level, text]) => ({
      t: "Heading",
      level,
      children: [{ t: "Text", value: text }],
    }),
  });
  registerFn.command({
    name: "bold",
    kind: "inline",
    schema: z.tuple([z.string()]),
    transform: (_n, [txt]) => ({
      t: "Bold",
      children: [{ t: "Text", value: txt }],
    }),
  });
}

export async function finalize(_ctx: PluginContext) {}
