import type { Op } from "../ops";
import { Renderer, type OpHandler } from "./renderer";

function textHandler(op: Op<any>): string {
  return `${op.data}`;
}

function fallbackHandler(op: Op<any>): string {
  return `unknown op: ${op.type}`;
}

export class TextRenderer extends Renderer {
  constructor() {
    super();
    this.addRenderFunc("text", textHandler);
  }

  override handleOp(op: Op<any>): string {
    const opHandler = (this.handlersMap.get(op.type) ?? fallbackHandler) as OpHandler<any>;
    return opHandler(op);
  }

  override render(ops: Op<any>[]): string {
    return ops.map(op => this.handleOp(op)).join("");
  }
}

export const textRenderer = new TextRenderer();
