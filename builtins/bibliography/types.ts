export type BibEntry = {
  id: string;
  type: string;
  fields: Record<string, string>;
};

export type BibState = {
  used: string[];
  bibliography: BibEntry[];
};
