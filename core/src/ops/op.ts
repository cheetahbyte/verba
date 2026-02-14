export type Op<T = unknown> = {
  type: string
  data: T
}

export type BaseOp = Op;

export function makeOp(type: string, data: any): BaseOp {
  return {type, data}
}
