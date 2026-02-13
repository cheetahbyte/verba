import { loadJSPlugin } from "./load-js";
import type { Host } from "./host";

export async function loadBuiltins(host: Host) {
  const hello = await import("@verba/builtin-hello");
  if (typeof (hello as any).register !== "function") {
    throw new Error("@verba/builtin-hello must export register(host)");
  }
  await (hello as any).register(host);
}
