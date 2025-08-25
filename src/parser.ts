import type { Node, TextNode, CommandNode, ParserError } from "./types";
import { ParserError as _ParserError } from "./types";

export function parse(input: string): Node[] {
  const p = new Scanner(input);
  const nodes: Node[] = [];

  while (!p.eof()) {
    const ch = p.peek();
    if (ch === ":" && p.peek(1) === ":" && !p.isEscaped()) {
      nodes.push(parseCommand(p));
    } else {
      nodes.push(parseText(p));
    }
  }

  return mergeAdjacentText(nodes);
}

// --- Scanner / Utils ---
class Scanner {
  constructor(
    public src: string,
    public i = 0,
  ) {}
  eof() {
    return this.i >= this.src.length;
  }
  peek(ahead = 0) {
    return this.src[this.i + ahead];
  }
  next() {
    return this.src[this.i++];
  }
  slice(start: number, end: number) {
    return this.src.slice(start, end);
  }
  isEscaped(idx = this.i): boolean {
    // Ein Token gilt als escaped, wenn eine ungerade Anzahl Backslashes davor steht
    let backslashes = 0;
    let j = idx - 1;
    while (j >= 0 && this.src[j] === "\\") {
      backslashes++;
      j--;
    }
    return backslashes % 2 === 1;
  }
  expect(char: string) {
    const got = this.next();
    if (got !== char)
      throw new _ParserError(
        `Expected '${char}', got '${got ?? "EOF"}'`,
        this.i - 1,
      );
  }
}

function parseText(p: Scanner): TextNode {
  const start = p.i;
  let out = "";
  while (!p.eof()) {
    if (p.peek() === ":" && p.peek(1) === ":" && !p.isEscaped()) break;
    out += p.next();
  }
  return { kind: "text", value: unescapeText(out), loc: { start, end: p.i } };
}

function parseCommand(p: Scanner): CommandNode {
  const start = p.i;
  // '::'
  p.expect(":");
  p.expect(":");

  const name = readIdent(p);
  if (!name) throw new _ParserError("Missing command name after '::'", p.i);

  skipSpaces(p);

  let args: string[] = [];
  if (p.peek() === "{") {
    const inner = readBalanced(p, "{", "}");
    args = splitArgs(inner).map((s) => s.trim());
  }

  const end = p.i;
  const raw = p.slice(start, end);
  return { kind: "command", name, args, raw, loc: { start, end } };
}

function readIdent(p: Scanner): string {
  let s = "";
  const isStart = (c?: string) => !!c && /[A-Za-z]/.test(c);
  const isPart = (c?: string) => !!c && /[A-Za-z0-9_-]/.test(c);
  if (!isStart(p.peek())) return "";
  s += p.next();
  while (isPart(p.peek())) s += p.next();
  return s;
}

function skipSpaces(p: Scanner) {
  while (!p.eof() && /\s/.test(p.peek())) p.next();
}

function readBalanced(p: Scanner, open: string, close: string): string {
  // Erwartet, dass der Cursor aktuell auf 'open' steht
  p.expect(open);
  let depth = 1;
  let out = "";
  let inStr: '"' | "'" | null = null;

  while (!p.eof()) {
    const ch = p.next();

    // String-Literale berücksichtigen ("..." / '...')
    if (!p.isEscaped(p.i - 1) && (ch === '"' || ch === "'")) {
      if (inStr === null) inStr = ch as '"' | "'";
      else if (inStr === ch) inStr = null;
      out += ch;
      continue;
    }

    if (inStr) {
      out += ch;
      continue;
    }

    if (ch === open && !p.isEscaped(p.i - 1)) {
      depth++;
      out += ch;
      continue;
    }
    if (ch === close && !p.isEscaped(p.i - 1)) {
      depth--;
      if (depth === 0) break;
      out += ch;
      continue;
    }

    out += ch;
  }

  if (depth !== 0) throw new _ParserError(`Unclosed '${open}'`, p.i);
  return out;
}

function splitArgs(s: string): string[] {
  const res: string[] = [];
  let buf = "";
  let depth = 0;
  let inStr: '"' | "'" | null = null;

  for (let i = 0; i < s.length; i++) {
    const ch = s[i];
    const prev = s[i - 1];
    const escaped = prev === "\\" && i - 1 >= 0;

    if (!escaped && (ch === '"' || ch === "'")) {
      if (inStr === null) inStr = ch as '"' | "'";
      else if (inStr === ch) inStr = null;
      buf += ch;
      continue;
    }

    if (inStr) {
      buf += ch;
      continue;
    }

    if (!escaped && ch === "{") {
      depth++;
      buf += ch;
      continue;
    }
    if (!escaped && ch === "}") {
      depth--;
      buf += ch;
      continue;
    }

    if (!escaped && depth === 0 && ch === ",") {
      res.push(buf);
      buf = "";
      continue;
    }

    buf += ch;
  }

  if (buf.length) res.push(buf);
  return res;
}

function unescapeText(s: string): string {
  // \:: -> ::, \{ -> {, \} -> }, \\ -> \  (nur das Nötigste für jetzt)
  return s
    .replace(/\\::/g, "::")
    .replace(/\\\{/g, "{")
    .replace(/\\\}/g, "}")
    .replace(/\\\\/g, "\\");
}

function mergeAdjacentText(nodes: Node[]): Node[] {
  const out: Node[] = [];
  for (const n of nodes) {
    const last = out[out.length - 1];
    if (last && last.kind === "text" && n.kind === "text") {
      (last as TextNode).value += n.value;
      (last as TextNode).loc.end = n.loc.end;
    } else out.push(n);
  }
  return out;
}
