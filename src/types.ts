export class ParserError extends Error {
  constructor(
    message: string,
    public pos: number,
  ) {
    super(`${message} @${pos}`);
    this.name = "ParserError";
  }
}

export type Position = number;

export type SourceLocation = {
  start: Position;
  end: Position;
};

export enum NodeType {
  Root = "root",
  Text = "text",
  Command = "command",
  Argument = "argument",
}

export interface BaseNode {
  type: NodeType;
  loc: SourceLocation;
}

export interface TextNode extends BaseNode {
  type: NodeType.Text;
  value: string;
}

export interface CommandNode extends BaseNode {
  type: NodeType.Command;
  name: string;
  arguments: ArgumentNode[];
}

export interface ArgumentNode extends BaseNode {
  type: NodeType.Argument;
  children: Node[];
}

export interface RootNode extends BaseNode {
  type: NodeType.Root;
  children: Node[];
}

export type Node = TextNode | CommandNode | ArgumentNode;
