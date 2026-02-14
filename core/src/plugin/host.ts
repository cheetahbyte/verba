import type { Context, OpHandler, PluginCommandHandler, Renderer, VerbaPlugin } from "..";
import { textRenderer } from "../render/text_renderer";

export class PluginHost {
  private commandMaps = new Map<string, PluginCommandHandler>();
  private renderersMap = new Map<string, Renderer>()

  registerCommand(
    name: string,
    handler: PluginCommandHandler,
    aliases?: string[]
  ): void {
    this.commandMaps.set(name, handler)
    aliases?.forEach(alias => this.commandMaps.set(alias, handler))
  }

  registerRenderer(name: string, renderer: Renderer) {
    this.renderersMap.set(name, renderer)
  }

  registerRenderHandlerFor<T>(
    rendererName: string,
    opName: string,
    func: OpHandler<T>
  ): void {
    const renderer = this.renderersMap.get(rendererName)
    if (!renderer) {
      console.error("failed to find renderer")
      return
    }
    renderer.addRenderFunc(opName, func)
  }

  getRenderer(name: string): Renderer {
    return this.renderersMap.get(name) ?? textRenderer
  }

  getSlice<T>(ctx: Context, pluginId: string, defaultState: T): T {
    if (!ctx.pluginData.has(pluginId)) {
      ctx.pluginData.set(pluginId, defaultState);
    }
    return ctx.pluginData.get(pluginId) as T;
  }

  async use(plugin: VerbaPlugin) {
    plugin.register?.(this);
  }

  async execute(ctx: Context, name: string, args: string[]) {
      const cmd = this.commandMaps.get(name);
      if (!cmd) throw new Error(`unknown command: ${name}`);
      return await cmd(ctx, args);
  }

  get commands(): string[] {
    return [...this.commandMaps.keys()]
  }

  get renderers(): string[] {
    return [...this.renderersMap.keys()]
}
}
