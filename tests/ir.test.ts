import { describe, it, expect } from "bun:test";
import { parse } from "@/parser";
import { blockize } from "@/blockize";
import { toIR } from "@/transform";
import { PluginLoader } from "@/plugins/loader";

it("produces paragraph + heading IR", async () => {
  const loader = new PluginLoader(1, { conflictPolicy: "warn" });
  await loader.loadAll();
  const ast = parse("::heading{1, Title}\n\nHello ::bold{world}");
  const { ir, diagnostics } = toIR(blockize(ast));
  expect(diagnostics).toEqual([]);
  expect(ir[0]).toMatchObject({ t: "Heading" });
  expect(ir[1]).toMatchObject({ t: "Paragraph" });
});
