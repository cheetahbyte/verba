export type Context = {
  pluginData: Map<string, any>;
}

export function newContext(): Context {
  return {
    pluginData: new Map()
  };
}
