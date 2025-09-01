import { Scanner } from "./scanner";
import { ArgumentNode, NodeType, RootNode, TextNode, type Node } from "./types";

export class Parser {
  constructor(private input: string) {}

  public parse(): RootNode {
    const scanner = new Scanner(this.input);
    const start = scanner.i;

    const children = this.parseUntil(scanner);

    return {
      type: NodeType.Root,
      children: children,
      loc: { start, end: scanner.i },
    };
  }

  private parseUntil(scanner: Scanner, endToken: string | null = null): Node[] {
    const nodes: Node[] = [];

    while (!scanner.eof()) {
      // stop when end token is found
      if (endToken && scanner.peek() === endToken) {
        break;
      }

      const nextChar = scanner.peek();
      if (nextChar === ":" && scanner.peek(1) === ":" && !scanner.isEscaped()) {
        nodes.push(this.parseCommand(scanner));
      } else {
        nodes.push(this.parseText(scanner, endToken));
      }
    }
    return nodes;
  }

  private parseText(scanner: Scanner, endToken: string | null): TextNode {
    const start = scanner.i;
    let out = "";
    while (!scanner.eof()) {
      if (
        (endToken && scanner.peek() === endToken) ||
        (scanner.peek() === ":" &&
          scanner.peek(1) === ":" &&
          !scanner.isEscaped())
      ) {
        break;
      }
      out += scanner.next();
    }
    return {
      type: NodeType.Text,
      value: out,
      loc: { start, end: scanner.i },
    };
  }

  parseCommand(scanner: Scanner): Node {
    const start = scanner.i;
    // consume the ::
    scanner.expect(":");
    scanner.expect(":");
    // parse name
    const nameStart = scanner.i;
    let name: string = "";
    while (!scanner.eof() && scanner.peek()?.match(/[a-zA-Z]/)) {
      name += scanner.next();
    }
    if (name.length == 0) {
      throw new Error("Expected command name after '::'");
    }

    // parse arguments
    const args: ArgumentNode[] = [];
    if (scanner.peek() === "{") {
      args.push(this.parseArguments(scanner));
    }

    return {
      type: NodeType.Command,
      arguments: args,
      loc: { start, end: scanner.i },
      name: name,
    };
  }

  parseArguments(scanner: Scanner): ArgumentNode {
    const start = scanner.i;
    scanner.expect("{");
    const children = this.parseUntil(scanner, "}");
    scanner.expect("}");
    return {
      type: NodeType.Argument,
      children: children,
      loc: { start, end: scanner.i },
    };
  }
}
