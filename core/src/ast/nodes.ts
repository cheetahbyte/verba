export type Span = { start: number; end: number };

export type DocumentNode = {
  kind: "document";
  span: Span;
  children: Node[];
};

export type TextNode = {
  kind: "text";
  span: Span;
  value: string;
};

export type CommandNode = {
  kind: "command";
  span: Span;
  name: string;
  args: ArgNode[];
};

export type ArgNode =
  | { kind: "argText"; span: Span; value: string };

export type Node = TextNode | CommandNode;
