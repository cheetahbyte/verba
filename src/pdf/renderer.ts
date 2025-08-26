import { PDFDocument, StandardFonts } from "pdf-lib";
import {
  LayoutDebug,
  LayoutLine,
  LayoutRun,
  PdfFonts,
  PdfOptions,
  RenderResult,
} from "./types";
import type { IRInline, IRNode } from "@/registry";
import { pageSizeToPoints, toPoints } from "./units";

function defaultOptions(): Required<PdfOptions> {
  return {
    pageSize: "A4",
    margins: { top: "20mm", right: "20mm", bottom: "25mm", left: "20mm" },
    fontSizes: { body: 11, h1: 24, h2: 18, h3: 14 },
    lineHeights: { body: 15, h1: 28, h2: 22, h3: 18 },
  };
}

export async function renderPdf(
  ir: IRNode[],
  opts: PdfOptions = {},
): Promise<RenderResult> {
  const pdfDoc = await PDFDocument.create();
  // Built-in fonts for Stage 1
  const bodyFont = await pdfDoc.embedFont(StandardFonts.Helvetica);
  const boldFont = await pdfDoc.embedFont(StandardFonts.HelveticaBold);
  const fonts: PdfFonts = { body: bodyFont, bold: boldFont };

  const options = { ...defaultOptions(), ...opts } as Required<PdfOptions>;
  const { width: pageW, height: pageH } = pageSizeToPoints(options.pageSize!);

  // Allow IR Margins to override provided options
  let marginTop = toPoints(options.margins!.top);
  let marginRight = toPoints(options.margins!.right);
  let marginBottom = toPoints(options.margins!.bottom);
  let marginLeft = toPoints(options.margins!.left);

  for (const node of ir) {
    if ((node as any).t === "Margins") {
      const m = node as any;
      marginLeft = toPoints(m.left);
      marginRight = toPoints(m.right);
      marginTop = toPoints(m.top);
      marginBottom = toPoints(m.bottom);
    }
  }

  const frameX = marginLeft;
  const frameYTop = pageH - marginTop;
  const frameWidth = pageW - marginLeft - marginRight;
  const frameBottom = marginBottom;

  let page = pdfDoc.addPage([pageW, pageH]);
  let cursorY = frameYTop;
  const debug: LayoutDebug = { pages: [{ lines: [] }] };

  const pushLine = (line: LayoutLine) => {
    const pageDebug = debug.pages[debug.pages.length - 1];
    pageDebug?.lines.push(line);
  };

  const ensureSpaceFor = (lineHeight: number) => {
    if (cursorY - lineHeight < frameBottom) {
      page = pdfDoc.addPage([pageW, pageH]);
      debug.pages.push({ lines: [] });
      cursorY = frameYTop;
    }
  };

  // Render each IR node
  for (const node of ir) {
    if ((node as any).t === "Margins") continue; // already applied

    if ((node as any).t === "Heading") {
      const h = node as any;
      const size =
        h.level === 1
          ? options.fontSizes!.h1
          : h.level === 2
            ? options.fontSizes!.h2
            : options.fontSizes!.h3;
      const lineHeight =
        h.level === 1
          ? options.lineHeights!.h1
          : h.level === 2
            ? options.lineHeights!.h2
            : options.lineHeights!.h3;
      const text = inlineCommandToPlainText(h.children);
      const lines = breakIntoLines(
        [{ text, fontKey: "bold", size, width: 0 }],
        frameWidth,
        fonts,
      );
      for (const line of lines) {
        ensureSpaceFor(lineHeight);
        drawLine(page, frameX, cursorY, line, fonts);
        pushLine({ ...line, lineHeight });
        cursorY -= lineHeight;
      }
      // spacing after heading
      cursorY -= Math.max(0, Math.round(lineHeight * 0.3));
      continue;
    }

    if ((node as any).t === "Paragraph") {
      const size = options.fontSizes!.body;
      const lineHeight = options.lineHeights!.body;
      const runs = inlineCommandToRuns((node as any).children, size);
      const lines = breakIntoLines(runs, frameWidth, fonts);
      for (const line of lines) {
        ensureSpaceFor(lineHeight);
        drawLine(page, frameX, cursorY, line, fonts);
        pushLine({ ...line, lineHeight });
        cursorY -= lineHeight;
      }
      // spacing after paragraph
      cursorY -= Math.max(0, Math.round(lineHeight * 0.2));
      continue;
    }

    // Unknown node types are ignored in Stage 1
  }

  const bytes = await pdfDoc.save();
  return { bytes, debug };
}

function inlineCommandToPlainText(children: IRInline[]): string {
  // TODO: more general implementation
  let out = "";
  for (const c of children) {
    if ((c as any).t === "Text") out += (c as any).value;
    else if ((c as any).t === "Bold")
      out += inlineCommandToPlainText((c as any).children);
  }
  return out;
}

function inlineCommandToRuns(children: IRInline[], size: number): LayoutRun[] {
  const runs: LayoutRun[] = [];
  for (const c of children) {
    if ((c as any).t === "Text") {
      runs.push({ text: (c as any).value, fontKey: "body", size, width: 0 });
    } else if ((c as any).t === "Bold") {
      runs.push({
        text: inlineCommandToPlainText((c as any).children),
        fontKey: "bold",
        size,
        width: 0,
      });
    }
  }
  return runs;
}

function tokenize(text: string): string[] {
  // Collapse any whitespace (space, tab, CR, LF) into a single space token
  const tokens: string[] = [];
  let buf = "";
  for (let i = 0; i < text.length; i++) {
    const ch = text[i];
    if (ch === " " || ch === "\t" || ch === "\n" || ch === "\r") {
      if (buf) {
        tokens.push(buf);
        buf = "";
      }
      if (tokens.length === 0 || tokens[tokens.length - 1] !== " ")
        tokens.push(" ");
    } else {
      buf += ch;
    }
  }
  if (buf) tokens.push(buf);
  return tokens;
}

function breakIntoLines(
  runs: LayoutRun[],
  maxWidth: number,
  fonts: PdfFonts,
): LayoutLine[] {
  const lines: LayoutLine[] = [];
  let curRuns: LayoutRun[] = [];
  let curWidth = 0;

  const pushLine = () => {
    lines.push({
      runs: curRuns,
      width: curWidth,
      lineHeight: estimateLineHeight(curRuns),
    });
    curRuns = [];
    curWidth = 0;
  };

  for (const run of runs) {
    const font = fonts[run.fontKey];
    const tokens = tokenize(run.text);
    for (const t of tokens) {
      const tokenWidth = font.widthOfTextAtSize(t, run.size);
      if (curWidth + tokenWidth <= maxWidth || curRuns.length === 0) {
        curRuns.push({ ...run, text: t, width: tokenWidth });
        curWidth += tokenWidth;
      } else {
        pushLine();
        curRuns.push({ ...run, text: t, width: tokenWidth });
        curWidth = tokenWidth;
      }
    }
  }
  if (curRuns.length) pushLine();
  return lines;
}

function estimateLineHeight(runs: LayoutRun[]): number {
  // simple: 1.35x of max size if mixed; otherwise font height could be used
  const maxSize = runs.reduce((m, r) => Math.max(m, r.size), 0);
  return Math.round(maxSize * 1.35);
}

function drawLine(
  page: import("pdf-lib").PDFPage,
  x: number,
  yTop: number,
  line: LayoutLine,
  fonts: PdfFonts,
) {
  let xPos = x;
  const baseline = yTop - Math.round(line.lineHeight * 0.8); // crude baseline; tweak later
  for (const r of line.runs) {
    const font = fonts[r.fontKey];
    if (!r.text) continue;
    page.drawText(r.text, { x: xPos, y: baseline, font, size: r.size });
    xPos += r.width;
  }
}
