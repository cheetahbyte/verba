import { Parser } from "./parser";
import { Scanner } from "./scanner";

async function main() {
  // get file content
  const file = Bun.file("example/sample.verba");
  const content = await file.text();
  // throw it into the parser
  const p = new Parser(content);
  const out = p.parse();
  console.log(out);
}

main();
