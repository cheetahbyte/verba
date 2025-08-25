export type Position = number;

export type SourceLocation = {
  start: Position;
  end: Position;
};

export enum NodeKind {
  Command = "command",
  Text = "text",
}

export type TextNode = {
  kind: NodeKind.Text;
  value: string;
  loc: SourceLocation;
};

export type CommandNode = {
  kind: NodeKind.Command;
  name: string;
  args: string[];
  raw: string; // original source code
  loc: SourceLocation;
};

export type Node = TextNode | CommandNode;

export class ParserError extends Error {
  constructor(
    message: string,
    public pos: number,
  ) {
    super(`${message} @${pos}`);
    this.name = "ParserError";
  }
}

export type LineCol = { line: number; col: number };

export function buildLineStarts(src: string): number[] {
  const starts = [0];
  for (let i = 0; i < src.length; i++) if (src[i] === "\n") starts.push(i + 1);
  return starts;
}
export function offsetToLineCol(offset: number, starts: number[]): LineCol {
  let lo = 0,
    hi = starts.length - 1;
  while (lo <= hi) {
    const mid = (lo + hi) >> 1;
    if (starts[mid]! <= offset) lo = mid + 1;
    else hi = mid - 1;
  }
  const line = hi;
  return { line: line + 1, col: offset - starts[line]! + 1 };
}

export type CommandOrigin = {
  name: string;
  version: string;
};
