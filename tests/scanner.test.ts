import { describe, it, expect } from "bun:test";
import { Scanner } from "@/scanner";
import { ParserError as _ParserError } from "@/types";

describe("Scanner", () => {
  it("eof/peek/next basics", () => {
    const s = new Scanner("abc");
    expect(s.eof()).toBe(false);

    expect(s.peek()).toBe("a");
    expect(s.peek(1)).toBe("b");
    expect(s.peek(2)).toBe("c");
    expect(s.peek(3)).toBeUndefined();

    expect(s.next()).toBe("a");
    expect(s.next()).toBe("b");
    expect(s.next()).toBe("c");
    expect(s.next()).toBeUndefined();

    expect(s.eof()).toBe(true);
  });

  it("slice returns substrings", () => {
    const s = new Scanner("Hello, Verba!");
    expect(s.slice(0, 5)).toBe("Hello");
    expect(s.slice(7, 12)).toBe("Verba");
  });

  it("isEscaped detects odd/even backslashes", () => {
    const s1 = new Scanner("\\::");
    const pos1 = s1.src.indexOf(":"); // first ':'
    s1.i = pos1;
    expect(s1.isEscaped()).toBe(true);

    // Two backslashes before ':' -> not escaped
    const s2 = new Scanner("\\\\::");
    const pos2 = s2.src.indexOf(":");
    s2.i = pos2;
    expect(s2.isEscaped()).toBe(false);

    const s3 = new Scanner("x{y}");
    const pos3 = s3.src.indexOf("{");
    s3.i = pos3;
    expect(s3.isEscaped()).toBe(false);

    const s4 = new Scanner("\\{foo}");
    const pos4 = s4.src.indexOf("{");
    s4.i = pos4;
    expect(s4.isEscaped()).toBe(true);
  });

  it("expect succeeds on matching char", () => {
    const s = new Scanner("::");
    expect(() => s.expect(":")).not.toThrow();
    expect(() => s.expect(":")).not.toThrow();
    expect(s.eof()).toBe(true);
  });

  it("expect throws ParserError on mismatch", () => {
    const s = new Scanner("ab");
    s.expect("a"); // consumes 'a'

    let err: unknown;
    try {
      s.expect("x"); // consumes 'b' und throws
    } catch (e) {
      err = e;
    }

    expect(err).toBeInstanceOf(_ParserError);
    expect(String(err)).toContain("Expected 'x', got 'b'");

    expect(s.i).toBe(2);
  });

  it("expect throws ParserError at EOF", () => {
    const s = new Scanner("");
    expect(() => s.expect(":")).toThrow();
    try {
      s.expect(":");
    } catch (e) {
      expect(e).toBeInstanceOf(_ParserError);
      expect(String(e)).toContain("EOF");
    }
  });
});
