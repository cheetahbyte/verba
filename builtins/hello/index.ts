import type { Host } from "@verba/core";

export function register(host: Host) {
  host.registerCommand("hello", (args) => {
    const name = args[0] ?? "world";
    return [`Hello, ${name}!`];
  });
}
