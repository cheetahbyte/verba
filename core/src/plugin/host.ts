export type CommandHandler = (
  args: string[]
) => Promise<string[]> | string[];

export class Host {
  private commands = new Map<string, CommandHandler>();

  registerCommand(name: string, handler: CommandHandler) {
    if (this.commands.has(name)) {
      throw new Error(`command already registered: ${name}`);
    }
    this.commands.set(name, handler);
  }

  async execute(name: string, args: string[]) {
    const cmd = this.commands.get(name);
    if (!cmd) throw new Error(`unknown command: ${name}`);
    return await cmd(args);
  }

  listCommands() {
    return [...this.commands.keys()];
  }
}
