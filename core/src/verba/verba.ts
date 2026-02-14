
import { newContext } from "../context";
import { evalDocumentOps } from "../eval";
import { parseDocument } from "../parser";
import { loadBuiltins } from "../plugin/discover";
import { PluginHost } from "../plugin/host";
import type { Renderer } from "../render";
import { textRenderer, TextRenderer } from "../render/text_renderer";

export class Verba {
  private pluginHost = new PluginHost()
  private currentRenderer: string = "text"
  private initialized = false;

  async init() {
      if (this.initialized) return;
      await loadBuiltins(this.pluginHost);
      this.initialized = true;
    }

  setRenderer(renderer: string) {
    this.currentRenderer = renderer
  }

  get renderer(): Renderer {
    return this.pluginHost.getRenderer(this.currentRenderer)
  }

  // executes the pipeline
  async execute(content: string) {
    // make sure plugins are loaded
    await this.init();
    // pipeline
    const documentContext = newContext()
    const ast = parseDocument(content);
    const ops = await evalDocumentOps(this.pluginHost, ast, documentContext)
    return this.renderer.render(ops)
  }
}
