import type { Node, TextNode, CommandNode } from "@/types";
import { NodeKind } from "@/types";
import { getSpec } from "@/registry";

export type Block =
  | { kind: "paragraph"; parts: (TextNode | CommandNode)[] }
  | { kind: "command"; node: CommandNode };

export function blockize(ast: Node[]): Block[] {
  const out: Block[] = [];
  let cur: (TextNode | CommandNode)[] = [];

  const flushPara = () => {
    if (!cur.length) return;
    const onlyWhitespace = cur.every(
      (n) => n.kind === NodeKind.Text && n.value.trim() === "",
    );
    if (!onlyWhitespace) out.push({ kind: "paragraph", parts: cur });
    cur = [];
  };

  for (const n of ast) {
    if (n.kind === NodeKind.Command) {
      const spec = getSpec(n.name);
      if (spec?.kind === "block") {
        flushPara();
        out.push({ kind: "command", node: n });
        continue;
      }
      cur.push(n);
      continue;
    }

    const chunks = n.value.split(/\n{2,}/);
    chunks.forEach((chunk, i) => {
      if (chunk) cur.push({ ...n, value: chunk });
      if (i < chunks.length - 1) flushPara();
    });
  }
  flushPara();
  return out;
}
