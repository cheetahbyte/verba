// scanner.test.ts
import { Scanner } from "@/scanner";
import { expect, test, describe } from "bun:test";

describe("Scanner Instantiation and Basic Reading", () => {
  test("scanner can be instantiated", () => {
    const testString = "hello world";
    const s = new Scanner(testString);
    expect(s.src).toBe(testString);
    expect(s.i).toBe(0);
  });

  test("peek() should look ahead without consuming", () => {
    const s = new Scanner("abc");
    expect(s.peek()).toBe("a");
    expect(s.i).toBe(0); // Position should not change
    expect(s.peek(1)).toBe("b");
    expect(s.i).toBe(0); // Position should still not change
  });

  test("next() should consume and return the next character", () => {
    const s = new Scanner("abc");
    expect(s.next()).toBe("a");
    expect(s.i).toBe(1);
    expect(s.next()).toBe("b");
    expect(s.i).toBe(2);
  });
});

describe("Scanner End-of-File (EOF) Behavior", () => {
  test("eof() should be false when not at the end", () => {
    const s = new Scanner("a");
    expect(s.eof()).toBe(false);
  });

  test("eof() should be true when at the end", () => {
    const s = new Scanner("a");
    s.next();
    expect(s.eof()).toBe(true);
  });

  test("peek() should return undefined at EOF", () => {
    const s = new Scanner("a");
    s.next();
    expect(s.peek()).toBe(undefined);
    expect(s.peek(5)).toBe(undefined); // Peeking past EOF
  });

  test("next() should return undefined at EOF", () => {
    const s = new Scanner("a");
    s.next(); // Consume 'a'
    expect(s.next()).toBe(undefined); // Now at EOF
    expect(s.i).toBe(2); // Index increases past length
  });

  test("should handle empty strings correctly", () => {
    const s = new Scanner("");
    expect(s.eof()).toBe(true);
    expect(s.peek()).toBe(undefined);
    expect(s.next()).toBe(undefined);
  });
});

describe("isEscaped() Method", () => {
  test("should return true for a single preceding backslash", () => {
    const s = new Scanner("a\\:b");
    expect(s.isEscaped(2)).toBe(true);
  });

  test("should return false for two preceding backslashes", () => {
    const s = new Scanner("a\\\\:b");
    expect(s.isEscaped(3)).toBe(false);
  });

  test("should return true for an odd number of preceding backslashes", () => {
    const s = new Scanner("a\\\\\\:b");
    expect(s.isEscaped(4)).toBe(true);
  });

  test("should return false if not preceded by a backslash", () => {
    const s = new Scanner("a:b");
    expect(s.isEscaped(1)).toBe(false);
  });

  test("should return false at the beginning of the string", () => {
    const s = new Scanner(":ab");
    expect(s.isEscaped(0)).toBe(false);
  });
});

describe("expect() Method", () => {
  test("should do nothing and advance if the character matches", () => {
    const s = new Scanner("abc");
    s.expect("a");
    expect(s.peek()).toBe("b");
    expect(s.i).toBe(1);
  });

  test("should throw an error if the character does not match", () => {
    const s = new Scanner("abc"); // Bun's syntax for testing thrown errors
    expect(() => s.expect("x")).toThrow("Expected 'x', got 'a'");
  });

  test("should throw an error when expecting at EOF", () => {
    const s = new Scanner("");
    expect(() => s.expect("a")).toThrow("Expected 'a', got 'EOF'");
  });
});
