import { z } from "zod";
import type { CommandNode } from "@/types";

export type IRNode =
  | { t: "Paragraph"; children: IRInline[] }
  | { t: "Heading"; level: number; children: IRInline[] }
  | { t: "Margins"; left: string; right: string; top: string; bottom: string }
  | { t: "BibliographyLoad"; file: string }
  | IRInline;

export type IRInline =
  | { t: "Text"; value: string }
  | { t: "Bold"; children: IRInline[] }
  | { t: "Cite"; value: string };

export type CommandKind = "block" | "inline";

export type CommandSpec<TSchema extends z.ZodTypeAny = z.ZodTypeAny> = {
  name: string;
  kind: CommandKind;
  schema: TSchema;
  transform: (node: CommandNode, parsed: z.infer<TSchema>) => IRNode;
};

const REG = new Map<string, CommandSpec>();

export function register<TSchema extends z.ZodTypeAny>(
  spec: CommandSpec<TSchema>,
) {
  REG.set(spec.name, spec as CommandSpec);
}

export function getSpec(name: string) {
  return REG.get(name);
}
register({
  name: "heading",
  kind: "block",
  schema: z.tuple([z.number().int().min(1).max(6), z.string()]),
  transform: (_n, [level, text]) => ({
    t: "Heading",
    level,
    children: [{ t: "Text", value: text }],
  }),
});

register({
  name: "margin",
  kind: "block",
  schema: z.tuple([z.string(), z.string(), z.string(), z.string()]),
  transform: (_n, [l, r, t, b]) => ({
    t: "Margins",
    left: l,
    right: r,
    top: t,
    bottom: b,
  }),
});

register({
  name: "bold",
  kind: "inline",
  schema: z.tuple([z.string()]),
  transform: (_n, [txt]) => ({
    t: "Bold",
    children: [{ t: "Text", value: txt }],
  }),
});
