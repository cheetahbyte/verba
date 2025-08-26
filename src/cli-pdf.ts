#!/usr/bin/env bun
import { parse } from "@/parser";
import { blockize } from "@/blockize";
import { toIR } from "@/transform";
import { PluginLoader } from "@/plugins/loader";
import { renderPdf } from "@/pdf/renderer";
import { writeFile } from "node:fs/promises";

async function main() {
  const [, , input, output = "output.pdf"] = Bun.argv;
  if (!input) {
    console.error("Usage: verba-pdf <input.verba> [output.pdf]");
    process.exit(1);
  }
  const src = await Bun.file(input).text();

  // Load core plugins (commands)
  const loader = new PluginLoader(1, { conflictPolicy: "warn" });
  await loader.loadAll();

  // Parse → IR
  const ast = parse(src);
  const blocks = blockize(ast);
  const { ir, diagnostics } = toIR(blocks);
  if (diagnostics.some((d) => d.level === "error")) {
    console.error("Diagnostics:");
    for (const d of diagnostics)
      console.error(`${d.level} ${d.code}: ${d.message}`);
  }

  // Render PDF
  const { bytes } = await renderPdf(ir);
  await writeFile(output, bytes);
  console.log(`✅ Wrote ${output}`);
}

main().catch((e) => {
  console.error(e);
  process.exit(1);
});
