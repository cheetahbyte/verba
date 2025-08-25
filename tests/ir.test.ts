import { describe, it, expect } from "bun:test";
import { parse } from "@/parser";
import { blockize } from "@/blockize";
import { toIR } from "@/transform";

it("produces paragraph + margins IR", () => {
  const ast = parse("::margin{20mm,20mm,25mm,25mm}\n\nHello ::bold{world}");
  const { ir, diagnostics } = toIR(blockize(ast));
  expect(diagnostics).toEqual([]);
  expect(ir[0]).toMatchObject({ t: "Margins" });
  expect(ir[1]).toMatchObject({ t: "Paragraph" });
});
