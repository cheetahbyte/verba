# Command Specification
A command starts with ::, followed by the command name and optional arguments in {…}.
Example (block + inline):
```
::heading{1, Introduction}

This is ::bold{important}.
```
## Registering Commands
Arguments are validated with Zod. The parser provides raw strings for arguments (comma-separated, already trimmed). If you need numbers (etc.), coerce in the schema.
#### Example: `heading` (block)
```ts
import { z } from "zod";

const heading = {
  name: "heading",
  kind: "block",
  schema: z.tuple([
    z.coerce.number().int().min(1).max(6), // "1" -> 1
    z.string(),                            // title text
  ]),
  transform: (_n, [level, text]) => ({
    t: "Heading",
    level,
    children: [{ t: "Text", value: text }],
  }),
};
```
#### Example: `block` (inline)
```ts
import { z } from "zod";

const bold = {
  name: "bold",
  kind: "inline",
  schema: z.tuple([z.string()]),
  transform: (_n, [txt]) => ({
    t: "Bold",
    children: [{ t: "Text", value: txt }],
  }),
};
```
> **Guidelines**
> - Inline commands must return an IRInline node (Bold, Cite, Text, …), not a Paragraph.
> - Block commands return top-level IR nodes (Heading, Margins, …).

### Command Kinds
- `inline`: used within paragraphs (formatting, short inlines like `bold`, `cite`, …)
- `block`: standalone blocks at the document level (`heading`, `margin`, `bibliography`, …)

### Arguments and Syntax
- **Form**: `::name{arg1, arg2, ...}`
- **Separation of Arguments**: commas at top level (nested {…} and quotes are respected)
- **Whitespace**: each argument is trimmed
- **Strings & Quotes**: both `'...'` and `"..."` are supported.
- **Escapes**: `\::`, `\{`, `\}` disable special meaning
- **Numbers**: arrive as strings; use `z.coerce.number()` in the command schema
- **Nesting**: inline commands can appear inside texty args (e.g., ::bold{see ::cite{"Doe 2020"}}).

### Namespacing (Plugin Commands)
In order to avoid collisions, plugin may namespace their command names
- **Syntax**: `::pluginName/commandName{}`
- **Core Commands**: Verbas core commands to not require a prefix
