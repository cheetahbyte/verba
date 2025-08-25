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
