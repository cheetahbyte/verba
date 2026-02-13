import type { DocumentNode, Node } from "../ast/nodes";
import { newContext, type Context } from "../context";
import type { Host } from "../plugin/host";

export async function evalDocument(host: Host, doc: DocumentNode): Promise<string> {
  let out = "";
  let documentContext: Context = newContext()

  for (const node of doc.children) {
    out += await evalNode(host, documentContext, node);
  }

  return out;
}

async function evalNode(host: Host, documentContext: Context, node: Node): Promise<string> {
  if (node.kind === "text") return node.value;

  const args = node.args.map((a) => a.value);
  const lines = await host.execute(documentContext, node.name, args);
  return lines.join("\n");
}
