import type { DocumentNode, PluginHost } from "..";
import { type Context, newContext } from "../context";
import type { BaseOp } from "../ops";
import type { Node } from "..";

export async function evalDocumentOps(
  host: PluginHost,
  doc: DocumentNode,
  documentContext: Context = newContext()
): Promise<BaseOp[]> {
  const ops: BaseOp[] = [];

  for (const node of doc.children) {
    ops.push(...(await evalNode(host, documentContext, node)));
  }

  return ops;
}

async function evalNode(
  host: PluginHost,
  documentContext: Context,
  node: Node
): Promise<BaseOp[]> {
  if (node.kind === "text") {
    return [{ type: "text", data: node.value }];
  }

  const args = node.args.map((a) => a.value);
  return await host.execute(documentContext, node.name, args);
}
