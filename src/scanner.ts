import { ParserError as _ParserError } from "./types";

export class Scanner {
  constructor(
    public src: string,
    public i = 0,
  ) {}

  eof() {
    return this.i >= this.src.length;
  }

  peak(ahead: number = 0) {
    return this.src[this.i + ahead];
  }

  next() {
    return this.src[this.i++];
  }

  slice(start: number, end: number) {
    return this.src.slice(start, end);
  }

  isEscaped(idx = this.i): boolean {
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
