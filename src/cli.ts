#!/usr/bin/env bun
import { parse } from "@/parser";

async function main() {
  const file = Bun.argv[2];
  if (!file) {
    console.error("Usage: verba-parse <file>");
    process.exit(1);
  }
  const src = await Bun.file(file).text();
  const ast = parse(src);
  console.log(JSON.stringify(ast, null, 2));
}

main().catch((err) => {
  console.error(err);
  process.exit(1);
});
