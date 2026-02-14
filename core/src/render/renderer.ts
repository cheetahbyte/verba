import type { Op } from "../ops";

export abstract class Renderer {
  constructor() {
    this.addRenderFunc = this.addRenderFunc.bind(this);
  }

  protected handlersMap = new Map<string, OpHandler<any, any>>();

  abstract render(ops: Op<any>[]): any;

  abstract handleOp(op: Op<any>): any;

  addRenderFunc<T>(opType: string, func: OpHandler<T>): void {
    this.handlersMap.set(opType, func);
  }

  get handlers(): string[] {
    return [...this.handlersMap.keys()]
}
}


export type OpHandler<T = unknown, R = any> = (op: Op<T>) => R;
