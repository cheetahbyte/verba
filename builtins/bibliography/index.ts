import type { Host } from "@verba/core";
import { readFile } from "node:fs/promises";
import Cite from "citation-js";
import type { Context } from "@verba/core/src/context";

type BibEntry = {
  id: string;
  type: string;
  fields: Record<string, string>;
};

type BibState = {
  used: string[];
  bibliography: BibEntry[];
};

export class BibliographyPlugin {
  private host: Host;

  constructor(host: Host) {
    this.host = host;
  }

  public register() {
    this.host.registerCommand("cite", (ctx, args) => this.handleCite(ctx, args), ["c"]);
    this.host.registerCommand("biblio", (ctx, args) => this.handleBiblio(ctx, args), ["bib"]);
    this.host.registerCommand("includebib", (ctx, args) => this.handleIncludeBib(ctx, args));
  }

  private getState(ctx: Context): BibState {
      return this.host.getSlice(ctx, "biblio", {
        used: [],
        bibliography: [],
      });
    }

    private handleCite(ctx: Context, args: string[]) {
      const state = this.getState(ctx);
      args.forEach((arg) => state.used.push(arg));
      return ["quoted " + args.join(",")];
    }

    private handleBiblio(ctx: Context, _args: string[]) {
      const state = this.getState(ctx);
      const uniqueUsedIds = [...new Set(state.used)];

      const outputLines = uniqueUsedIds.map((id, index) => {
        const entry = state.bibliography.find((e) => e.id === id);
        if (!entry) return `[${index + 1}] Key "${id}" not found.`;

        const { author = "Unknown Author", title = "No Title", year = "n.d." } = entry.fields;
        return `[${index + 1}] ${author} (${year}). ${title}.`;
      });

      return outputLines.length > 0 ? ["--- BIBLIOGRAPHY ---", ...outputLines] : [];
    }

    private async handleIncludeBib(ctx: Context, args: string[]) {
      const state = this.getState(ctx);
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
          },
        }));

        newEntries.forEach((entry) => {
          if (!state.bibliography.find((e) => e.id === entry.id)) {
            state.bibliography.push(entry);
          }
        });

        return [`Loaded ${newEntries.length} entries via citation.js`];
      } catch (err) {
        return [`Error: ${(err as Error).message}`];
      }
    }
  }

export function register(host: Host) {
  const plugin = new BibliographyPlugin(host);
  plugin.register();
}
