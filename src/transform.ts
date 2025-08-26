import type { IRNode, IRInline } from "@/registry";
import type { Block } from "@/blockize";
import { getSpec } from "@/registry";
import { z } from "zod";
import { Diagnostic, Diagnostics } from "@/diagnostics";

export type TransformResult = { ir: IRNode[]; diagnostics: Diagnostic[] };

export function toIR(
  blocks: Block[],
  opts?: { reporter?: Diagnostics; lineStarts?: number },
): TransformResult {
  const ir: IRNode[] = [];
  const diagnostics: Diagnostic[] = [];
  const report = (d: Diagnostic) => {
    diagnostics.push(d);
    opts?.reporter?.report?.(d);
  };
  for (const b of blocks) {
    if (b.kind === "command") {
      const spec = getSpec(b.node.name);
      if (!spec) {
        report({
          level: "error",
          code: "transform/unknown-command",
          message: `Unknown command ::${b.node.name}`,
          pos: b.node.loc?.start,
        });
        continue;
      }
      try {
        const parsed = (spec.schema as z.ZodTypeAny).parse(b.node.args);
        ir.push(spec.transform(b.node, parsed));
      } catch (e) {
        report({
          level: "error",
          code: "transform/invalid-args",
          message: `::${b.node.name} arguments invalid: ${(e as Error).message}`,
          plugin: spec.origin
            ? `${spec.origin.name}@${spec.origin.version}`
            : undefined,
          pos: b.node.loc?.start,
        });
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
          report({
            level: "error",
            code: "transform/unknown-inline",
            message: `Unknown inline ::${part.name}`,
            pos: part.loc?.start,
          });
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
              report({
                level: "warn",
                code: "transform/block-used-inline",
                message: `Block command ::${part.name} used inline`,
                plugin: spec.origin
                  ? `${spec.origin.name}@${spec.origin.version}`
                  : undefined,
                pos: part.loc?.start,
              });
            }
          }
        } catch (e) {
          report({
            level: "error",
            code: "transform/invalid-args",
            message: `::${part.name} arguments invalid: ${(e as Error).message}`,
            plugin: spec.origin
              ? `${spec.origin.name}@${spec.origin.version}`
              : undefined,
            pos: part.loc?.start,
          });
        }
      }
    }
    ir.push({ t: "Paragraph", children });
  }

  return { ir, diagnostics };
}
