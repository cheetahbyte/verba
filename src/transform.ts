// src/transform.ts
import type { IRNode, IRInline } from "@/registry";
import type { Block } from "@/blockize";
import { getSpec } from "@/registry";
import { ParserError } from "@/types";
import { z } from "zod";

export type TransformResult = { ir: IRNode[]; diagnostics: string[] };

export function toIR(blocks: Block[]): TransformResult {
  const ir: IRNode[] = [];
  const diagnostics: string[] = [];

  for (const b of blocks) {
    if (b.kind === "command") {
      const spec = getSpec(b.node.name);
      if (!spec) {
        diagnostics.push(`Unknown command ::${b.node.name}`);
        continue;
      }
      try {
        const parsed = (spec.schema as z.ZodTypeAny).parse(b.node.args);
        ir.push(spec.transform(b.node, parsed));
      } catch (e) {
        diagnostics.push(
          `Invalid args for ::${b.node.name}: ${(e as Error).message}`,
        );
      }
      continue;
    }

    // paragraph
    const children: IRInline[] = [];
    for (const part of b.parts) {
      if (part.kind === "text") {
        children.push({ t: "Text", value: part.value });
      } else {
        const spec = getSpec(part.name);
        if (!spec) {
          diagnostics.push(`Unknown inline ::${part.name}`);
          continue;
        }
        try {
          const parsed = (spec.schema as z.ZodTypeAny).parse(part.args);
          const node = spec.transform(part, parsed);
          // nur Inline zulassen
          if ("t" in node && typeof node.t === "string") {
            if (node.t === "Bold" || node.t === "Cite" || node.t === "Text") {
              children.push(node as IRInline);
            } else {
              diagnostics.push(`Block command ::${part.name} used inline`);
            }
          }
        } catch (e) {
          diagnostics.push(
            `Invalid args for ::${part.name}: ${(e as Error).message}`,
          );
        }
      }
    }
    ir.push({ t: "Paragraph", children });
  }

  return { ir, diagnostics };
}
