import type { BaseOp, Context, VerbaPlugin, VerbaPluginHost } from "@verba/core";
import { makeOp } from "@verba/core";
import { readFile } from "node:fs/promises";
import Cite from "citation-js";
import type { BibState, BibEntry } from "./types";
import { renderOpCiteText } from "./renderfuncs";

const PLUGIN_ID = "biblio";

function getState(host: VerbaPluginHost, ctx: Context): BibState {
  return host.getSlice(ctx, PLUGIN_ID, {
    used: [],
    bibliography: [],
  });
}

function handleCite(host: VerbaPluginHost, ctx: Context, args: string[]): BaseOp[] {
  const state = getState(host, ctx);
  args.forEach((arg) => state.used.push(arg));
  const quoted = state.bibliography.filter((entry) =>
     args.includes(entry.id)
   );
  return [makeOp("citation", {quoted})];
}

function handleBiblio(host: VerbaPluginHost, ctx: Context): BaseOp[] {
  const state = getState(host, ctx);
  const uniqueUsedIds = [...new Set(state.used)];
  const entries = uniqueUsedIds.map((id, index) => state.bibliography.find((e) => e.id === id));

  return [makeOp("bibliography", {used: entries})];
}

async function handleIncludeBib(
  host: VerbaPluginHost,
  ctx: Context,
  args: string[]
): Promise<BaseOp[]> {
  const state = getState(host, ctx);
  const filePath = args[0];
  if (!filePath) return [makeOp("text", "Error: No bib file specified")];

  try {
    const content = await readFile(filePath, "utf-8");
    const data = new Cite(content);

    const newEntries: BibEntry[] = data.data.map((entry: any) => ({
      id: entry.id,
      type: entry.type,
      fields: {
        author: entry.author?.map((a: any) => `${a.given} ${a.family}`).join(", "),
        title: entry.title,
        year: entry.issued?.["date-parts"]?.[0]?.[0]?.toString(),
        journal: entry["container-title"],
        publisher: entry.publisher,
      },
    }));

    newEntries.forEach((entry) => {
      if (!state.bibliography.find((e) => e.id === entry.id)) {
        state.bibliography.push(entry);
      }
    });
    return []
  } catch (err) {
    console.error("file not found")
    return []
  }
}

export const bibliographyPlugin: VerbaPlugin = {
  id: "builtin-bibliography",
  register(host) {
    host.registerCommand("cite", (ctx, args) => handleCite(host, ctx, args), ["c"]);
    host.registerCommand("biblio", (ctx, _args) => handleBiblio(host, ctx), ["bib"]);
    host.registerCommand("includebib", (ctx, args) => handleIncludeBib(host, ctx, args));
    // render handlers
    host.registerRenderHandlerFor("text", "citation", renderOpCiteText)
  },

};

export default bibliographyPlugin;
