import type { Op } from "@verba/core";
import type { BibEntry } from "./types";


export type CiteOpPayload = {
  quoted: BibEntry[]
}

export function renderOpCiteText(op: Op<CiteOpPayload>) {
  return `(${op.data.quoted.at(0)!.fields.author}, ${op.data.quoted.at(0)!.fields.year})`
}
