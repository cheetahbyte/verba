import { loadJSPlugin } from "./load-js";
import type { Host } from "./host";

export async function loadBuiltins(host: Host) {
  // bibliography
  const biblio = await import("@verba/builtin-bibliography")
  if (typeof (biblio as any).register !== "function") {
    throw new Error("@verba/builtin-hello must export register(host)");
  }
  await (biblio as any).register(host);
}
