import type { ArgNode, CommandNode, DocumentNode, Node, Span } from "../ast";

function span(start: number, end: number): Span {
  return { start, end };
}

function isIdentChar(ch: string) {
  return /[a-zA-Z0-9_-]/.test(ch);
}

function parseArgs(raw: string, baseOffset: number): ArgNode[] {
  // minimal: split by commas, trim whitespace
  // later you can handle quoting/escaping/nesting
  const out: ArgNode[] = [];
  let cursor = 0;

  for (const part of raw.split(",")) {
    const start = cursor;
    cursor += part.length + 1; // + comma
    const value = part.trim();
    const pStart = baseOffset + start;
    const pEnd = baseOffset + start + part.length;

    if (value.length === 0) continue;
    out.push({ kind: "argText", span: span(pStart, pEnd), value });
  }
  return out;
}

export function parseDocument(source: string): DocumentNode {
  const nodes: Node[] = [];

  let i = 0;
  let textStart = 0;

  const flushText = (end: number) => {
    if (end <= textStart) return;
    const value = source.slice(textStart, end);
    if (value.length === 0) return;
    nodes.push({ kind: "text", span: span(textStart, end), value });
  };

  while (i < source.length) {
    // command begins with ::
    if (source[i] === ":" && source[i + 1] === ":") {
      flushText(i);

      const cmdStart = i;
      i += 2;

      // parse name
      const nameStart = i;
      while (i < source.length && isIdentChar(source[i])) i++;
      const name = source.slice(nameStart, i);

      // require { ... }
      if (source[i] !== "{") {
        // treat as text if malformed
        // rollback to include ::
        textStart = cmdStart;
        i = cmdStart + 2;
        continue;
      }

      i++; // skip '{'
      const argsStart = i;

      // find matching '}' (no nesting yet)
      while (i < source.length && source[i] !== "}") i++;
      const argsEnd = i;

      if (source[i] !== "}") {
        // unterminated => treat as text
        textStart = cmdStart;
        break;
      }
      i++; // skip '}'

      const rawArgs = source.slice(argsStart, argsEnd);
      const args = parseArgs(rawArgs, argsStart);

      const cmd: CommandNode = {
        kind: "command",
        span: span(cmdStart, i),
        name,
        args
      };

      nodes.push(cmd);
      textStart = i;
      continue;
    }

    i++;
  }

  flushText(source.length);

  return {
    kind: "document",
    span: span(0, source.length),
    children: nodes
  };
}
