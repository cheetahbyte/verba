export type UnitString = `${number}${"mm" | "cm" | "in" | "pt"}`;

const PT_PER_IN = 72;
const MM_PER_IN = 25.4;

// points are used since they provide a consistent measure for layouts. further reading: https://www.theinformedillustrator.com/2016/08/point-pixel-pica-simply-defined.html#:~:text=A%20“pica”%20is%20a%20unit,text%20width%2C%20spaces%2C%20etc.&text=Highlighted%20area%20shows%20the%20number%20of%20picas%20in%20an%208.5,use%20of%20fractions%20or%20decimals.
export function toPoints(input: number | string): number {
  // we already have points
  if (typeof input == "number") return input;
  const trimmed = input.trim();
  const match = trimmed.match(/^(-?\d+(?:\.\d+)?)(mm|cm|in|pt)$/i);
  if (!match) throw new Error(`Invalid unit '${input}'. Use mm|cm|in|pt`);
  const value = Number(match[1]);
  const unit = match[2]?.toLowerCase();
  switch (unit) {
    case "pt":
      return value;
    case "in":
      return value * PT_PER_IN;
    case "mm":
      return value * (PT_PER_IN / MM_PER_IN);
    case "cm":
      return value * (PT_PER_IN / MM_PER_IN);
    default:
      throw new Error(`Unsupported unit '${unit}'`);
  }
}

export function pageSizeToPoints(
  size: "A4" | { widthPt: number; heightPt: number },
): { width: number; height: number } {
  if (typeof size === "string") {
    if (size === "A4") return { width: 595.28, height: 841.89 }; // ISO A4 in pts
  }
  return { width: (size as any).widthPt, height: (size as any).heightPt };
}

export function round(n: number, dp = 3) {
  const p = 10 ** dp;
  return Math.round(n * p) / p;
}
