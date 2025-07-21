# Requirements: Markdown to Backlog Converter CLI

## 1. Project Overview

Create a command-line interface (CLI) tool written in **Go** that converts Markdown text into Backlog notation. The tool must perform high-precision conversion by parsing the Markdown into an Abstract Syntax Tree (AST) rather than using simple string replacement.

This project is intended to be published as an Open Source Software on GitHub, so it must include high-quality code, documentation, and automation.

---

## 2. CLI Specifications

The tool should be invokable from the command line with the following specifications.

-   **Command Name:** `md2backlog`

-   **Functionality:**
    -   Read Markdown from a specified input file.
    -   Read Markdown from standard input (`stdin`) if no file is specified.
    -   Write the converted Backlog notation to a specified output file.
    -   Write the converted output to standard output (`stdout`) if no output file is specified.

-   **Command-Line Flags:**
    -   `-i`, `--input <file>`: Specifies the input Markdown file path.
    -   `-o`, `--output <file>`: Specifies the output file path.
    -   `--help`: Displays the help message.
    -   `--version`: Displays the tool's version.

-   **Usage Examples:**
    ```bash
    # Convert file to standard output
    md2backlog -i input.md

    # Convert file to another file
    md2backlog -i input.md -o output.txt

    # Convert from standard input to standard output
    cat input.md | md2backlog
    ```

---

## 3. Conversion Rules

The following table details the required conversion rules from Markdown to Backlog notation.

| Category | Markdown | Backlog | Conversion Policy |
| :--- | :--- | :--- | :--- |
| **Headings** | `# H1`, `## H2`... | `* H1`, `** H2`... | Convert `#` count to `*` count. |
| **Emphasis** | `**Bold**` | `''Bold''` | |
| | `*Italic*` | `'''Italic'''` | |
| | `~~Strikethrough~~` | `%%Strikethrough%%` | |
| **Unordered Lists**| `- item` or `* item` | `- item` | Standardize to `-`. |
| | `- L1`<br>`  - L2`<br>`    - L3` | `- L1`<br>`-- L2`<br>`--- L3` | Use hyphen count for nesting level. |
| **Numbered Lists**| `1. item` | `+ item` | |
| | `1. L1`<br>`   1. L2` | `+ L1`<br`+ L2` | **Flatten nested lists** as Backlog does not support them. |
| **Code** | `` `inline code` `` | `{code}inline code{/code}` | |
| | \`\`\`lang<br>code<br>\`\`\` | `>{code:lang}<br>code<br>{/code}<` | |
| **Others** | `> Quote` | `> Quote` | No change. |
| | `[Link](URL)` | `[[Link:URL]]` | |
| | `![ALT](URL)` | `![ALT](URL)` | No change. |
| | `---` | `---` | No change. |
| | `|th1|...`<br>`|---|...`<br>`|td1|...` | `|*th1|...`<br`|td1|...` | Add `*` to header cells. |

---

## 4. Technical Stack & Repository Requirements

-   **Language:** **Go**
-   **Markdown Parser:** **goldmark**. Use its AST walking capabilities (`ast.Walk`) for the conversion logic.
-   **CLI Framework:** **Cobra**. Use it to build the command structure and handle flags.

-   **Code Quality & Formatting:**
    -   **Formatter:** All code must be formatted with `goimports`.
    -   **Linter:** **golangci-lint** must be used. A configuration file (`.golangci.yml`) should be included to define the linting rules.

-   **Testing:**
    -   Implement comprehensive **unit tests** for each conversion rule using the standard `testing` package.
    -   Implement **integration tests** that run the compiled CLI tool with sample files and verify the output.

-   **Automation (CI/CD):**
    -   Set up a **GitHub Actions** workflow for Continuous Integration (`.github/workflows/ci.yml`).
    -   The CI workflow must run on every push and pull request.
    -   The CI pipeline must execute: `gofmt` check, `golangci-lint`, and `go test`.
    -   Set up an automated release process using **GoReleaser**. This should be triggered by pushing a git tag (e.g., `v1.0.0`) and should create cross-compiled binaries for Windows, macOS, and Linux, attaching them to a GitHub Release.

-   **Documentation:**
    -   `README.md`: Must include project overview, installation instructions, usage, and the conversion rules table.
    -   `LICENSE`: Must include a permissive OSS license (e.g., MIT or Apache 2.0).
    -   `.gitignore`: Must be configured for a Go project.