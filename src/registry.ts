import { z } from "zod";
import type { CommandNode, CommandOrigin } from "@/types";

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
  origin: CommandOrigin;
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
