import type { Host } from "@verba/core";
import { readFile } from "node:fs/promises";
import Cite from "citation-js";

type BibEntry = {
  id: string;
  type: string;
  fields: Record<string, string>;
}

export function register(host: Host) {
  host.registerCommand("cite", (ctx, args) => {
    const state = host.getSlice(ctx, "biblio", { used: [] as string[], bibliography: [] as BibEntry[] });
    args.forEach(arg => state.used.push(arg))
    return ["quoted " + args.join(",")]
  }, ["c"]);

  host.registerCommand("biblio", (ctx, args) => {
    const state = host.getSlice(ctx, "biblio", {
      used: [] as string[],
      bibliography: [] as BibEntry[]
    });

    const uniqueUsedIds = [...new Set(state.used)];

    const outputLines = uniqueUsedIds.map((id, index) => {
      const entry = state.bibliography.find(e => e.id === id);

      if (!entry) {
        return `[${index + 1}] Key "${id}" not found in bibliography.`;
      }

      const author = entry.fields["author"] || "Unknown Author";
      const title = entry.fields["title"] || "No Title";
      const year = entry.fields["year"] || "n.d.";

      return `[${index + 1}] ${author} (${year}). ${title}.`;
    });

    if (outputLines.length === 0) {
      return [];
    }

    return [
      "--- BIBLIOGRAPHY ---",
      ...outputLines
    ];
  }, ["bib"]);

  host.registerCommand("includebib", async (ctx, args) => {
    const state = host.getSlice(ctx, "biblio", {
      used: [] as string[],
      bibliography: [] as BibEntry[]
    });

    const filePath = args[0];
    if (!filePath) return ["Error: No bib file specified"];

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
        }
      }));

      newEntries.forEach(entry => {
        if (!state.bibliography.find(e => e.id === entry.id)) {
          state.bibliography.push(entry);
        }
      });

      return [`Loaded ${newEntries.length} entries via citation.js`];
    } catch (err) {
      return [`Error: ${(err as Error).message}`];
    }
  });
}
