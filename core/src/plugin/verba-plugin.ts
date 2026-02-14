import type { Context } from "../context";
import type { BaseOp} from "../ops";
import type { OpHandler, Renderer } from "../render";

export type PluginCommandHandler = (
  ctx: Context,
  args: string[]
) => Promise<BaseOp[]> | BaseOp[];

export interface VerbaPluginHost {
  registerCommand(
    name: string,
    handler: PluginCommandHandler,
    aliases?: string[]
  ): void;

  registerRenderHandlerFor<T>(
    rendererName: string,
    opName: string,
    func: OpHandler<T>
  ): void;

  registerRenderer(name: string, renderer: Renderer): void;

  getSlice<T>(ctx: Context, pluginId: string, defaultState: T): T;
}

export interface VerbaPlugin {
  id: string;
  register?(host: VerbaPluginHost): void | Promise<void>;
}
