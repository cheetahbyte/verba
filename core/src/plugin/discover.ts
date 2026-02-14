import { textRenderer, TextRenderer } from "../render/text_renderer";
import type { PluginHost } from "./host";
import type { VerbaPlugin } from "./verba-plugin";

function isVerbaPlugin(value: unknown): value is VerbaPlugin {
  if (!value || typeof value !== "object") {
    return false;
  }

  const maybePlugin = value as Partial<VerbaPlugin>;
  return typeof maybePlugin.id === "string";
}

export async function loadBuiltins(host: PluginHost) {
  // renderers
  host.registerRenderer("text", textRenderer)

  // bibliography
  const biblio = await import("@verba/builtin-bibliography");

  const plugin = (biblio as any).default ?? (biblio as any).plugin;
  if (isVerbaPlugin(plugin)) {
    await host.use(plugin);
    return;
  }

  if (typeof (biblio as any).register !== "function") {
    throw new Error(
      "@verba/builtin-bibliography must export default/plugin VerbaPlugin or register(host)"
    );
  }
  await (biblio as any).register(host);
}
