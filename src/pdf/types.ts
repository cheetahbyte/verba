import type { IRNode, IRInline } from "@/registry";

export type PdfFonts = {
  body: import("pdf-lib").PDFFont;
  bold: import("pdf-lib").PDFFont;
};

export type PdfOptions = {
  pageSize?: "A4" | { widthPt: number; heightPt: number };
  margins?: { top: string; right: string; bottom: string; left: string };
  fontSizes?: { body: number; h1: number; h2: number; h3: number };
  lineHeights?: { body: number; h1: number; h2: number; h3: number };
};

export type LayoutRun = {
  text: string;
  fontKey: keyof PdfFonts;
  size: number;
  width: number;
};

export type LayoutLine = {
  runs: LayoutRun[];
  width: number;
  lineHeight: number;
};
export type LayoutPage = { lines: LayoutLine[] };
export type LayoutDebug = { pages: LayoutPage[] };
export type RenderResult = { bytes: Uint8Array; debug: LayoutDebug };
