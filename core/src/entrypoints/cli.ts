#!/usr/bin/env bun
import path from "node:path";
import { Verba } from "../verba/verba";


async function readStdin(): Promise<string> {
  const data = await new Response(process.stdin).text();
  return data;
}

async function main() {
  const verba = new Verba()

  const [, , file] = process.argv;

  const input = file
    ? await Bun.file(path.resolve(process.cwd(), file)).text()
    : await readStdin();
  const output = (await verba.execute(input)) as unknown as string

  process.stdout.write(output);
}

main().catch((err) => {
  console.error(err?.stack ?? String(err));
  process.exit(1);
});
