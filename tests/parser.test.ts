import { describe, expect, it } from "bun:test";
import { parse } from "@/parser";
import { NodeKind, type CommandNode, type TextNode } from "@/types";

describe("Verba Parser", () => {
  it("parses plain text", () => {
    const ast = parse("Hello World\n");
    expect(ast.length).toBe(1);
    const n = ast[0] as TextNode;
    expect(n.kind).toBe(NodeKind.Text);
    expect(n.value).toBe("Hello World\n");
  });

  it("parses single command with args", () => {
    const ast = parse("::command{Intro, v1}\n");
    const c = ast[0] as CommandNode;
    expect(c.kind).toBe(NodeKind.Command);
    expect(c.name).toBe("command");
    expect(c.args).toEqual(["Intro", "v1"]);
  });

  it("parses inline command between text", () => {
    const ast = parse("This is the ::bold{cool} command.\n");
    expect(ast.map((n) => n.kind)).toEqual([
      NodeKind.Text,
      NodeKind.Command,
      NodeKind.Text,
    ]);
  });

  it("supports nested braces in args", () => {
    const ast = parse("::wrap{::bold{deep}, level2}");
    const cmd = ast[0] as CommandNode;
    expect(cmd.args[0]).toBe("::bold{deep}");
  });

  it("ignores escaped markers", () => {
    const ast = parse("\\::not-a-cmd and real ::cmd{ok}");
    expect(ast.length).toBe(2);
    expect(ast[0]?.kind).toBe(NodeKind.Text);
    expect(ast[1]?.kind).toBe(NodeKind.Command);
  });
});
