#!/usr/bin/env bun
import { Host } from "../plugin/host";
import { loadBuiltins } from "../plugin/discover";
import { parseDocument } from "../parser";
import { evalDocument } from "../eval";
import path from "node:path";


async function readStdin(): Promise<string> {
  const data = await new Response(process.stdin).text();
  return data;
}

async function main() {
  const host = new Host();
  await loadBuiltins(host);

  const [, , file] = process.argv;

  const input = file
    ? await Bun.file(path.resolve(process.cwd(), file)).text()
    : await readStdin();

  const ast = parseDocument(input);
  const output = await evalDocument(host, ast);

  process.stdout.write(output);
}

main().catch((err) => {
  console.error(err?.stack ?? String(err));
  process.exit(1);
});
