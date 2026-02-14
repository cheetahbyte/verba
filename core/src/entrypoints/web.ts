import { Verba } from "../verba/verba";


export class VerbaWeb {
  private instance: Verba;

  constructor() {
    this.instance = new Verba();
  }

  async process(input: string): Promise<string> {
    try {
      const result = await this.instance.execute(input);
      return String(result);
    } catch (error) {
      throw new Error(`Verba Execution Failed: ${error instanceof Error ? error.message : String(error)}`);
    }
  }
}

export const createVerba = () => new VerbaWeb();
