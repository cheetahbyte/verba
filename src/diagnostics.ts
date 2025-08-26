import { Position } from "./types";

export type DiagnosticLevel = "info" | "warn" | "error";
export type Diagnostic = {
  level: DiagnosticLevel;
  code: string;
  message: string;
  plugin?: string;
  pos?: Position;
};

export interface Diagnostics {
  report(d: Diagnostic): void;
}

export function defaultDiagnostics(): Diagnostics {
  return {
    report(d) {
      const tag = d.level === "error" ? "❌" : d.level === "warn" ? "⚠️" : "ℹ️";
      console.error(
        `${tag} [${d.code}] ${d.message}${d.plugin ? ` (${d.plugin})` : ""}`,
      );
    },
  };
}
