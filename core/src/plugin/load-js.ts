import { pathToFileURL } from "node:url";
import type { Host } from "./host";

export async function loadJSPlugin(host: Host, absPath: string) {
  const mod = await import(pathToFileURL(absPath).href);

  if (typeof mod.register !== "function") {
    throw new Error("plugin must export register(host)");
  }

  await mod.register(host);
}
