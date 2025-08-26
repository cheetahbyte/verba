#!/usr/bin/env bun
import { parse } from "@/parser";
import { blockize } from "@/blockize";
import { toIR } from "@/transform";
import { PluginLoader } from "./plugins/loader";

const [, , cmd = "ast", file] = Bun.argv;
if (!file) {
  console.error("Usage: verba <ast|ir|lint> <file>");
  process.exit(1);
}
const src = await Bun.file(file).text();
const loader = new PluginLoader(1, { conflictPolicy: "warn" });
await loader.loadAll();
if (cmd === "ast") {
  console.log(JSON.stringify(parse(src), null, 2));
} else if (cmd === "ir") {
  const { ir, diagnostics } = toIR(blockize(parse(src)));
  if (diagnostics.length) console.error("Diagnostics:", diagnostics);
  console.log(JSON.stringify(ir, null, 2));
} else if (cmd === "lint") {
  const { diagnostics } = toIR(blockize(parse(src)));
  diagnostics.forEach((d) => console.error("•", d));
  process.exit(diagnostics.length ? 1 : 0);
}
