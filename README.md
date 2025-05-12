># Go SSG (Static Site Generator)

A simple static site generator built with Go. This project takes Markdown files from a `content` directory, converts them to HTML, and serves them from an `output` directory.

## Prerequisites

- Go (version 1.24.2 or later recommended)
- Make (optional, for using the Makefile)

## Configuration

Configuration can be managed via a YAML file and overridden by CLI flags for relevant commands.

Create a configuration file at `pkg/config/config.yaml` (or specify a different path using the `-c` flag).

Example `pkg/config/config.yaml`:
```yaml
siteTitle: "My Awesome Static Site"
contentDir: "content"
outputDir: "output"
serverPort: "8080"
```

**Configuration Options:**

-   `siteTitle`: The title of your site (currently used for logging, can be integrated into templates).
-   `contentDir`: The directory where your Markdown source files are located. Default: `content`.
-   `outputDir`: The directory where the generated HTML files will be saved. Default: `output`.
-   `serverPort`: The port on which the local development server will run. Default: `8080`.

## Setup

    ```bash
    git clone <your-repo-url>
    cd go-ssg
    ```

### Using Makefile

-   To build the application:
    ```bash
    make all 
    ```
    (Note: `make all` in your provided Makefile runs `clean`, `dirs`, `go mod tidy`, then `go build`)
-   To run a command (e.g., generate then serve):
    ```bash
    ./go-ssg generate
    ./go-ssg serve
    ```
    Or, to get help:
    ```bash
    make run 
    ```

-   To clean build artifacts and the `output` directory (content directory isn't removed):
    ```bash
    make clean
    ```

### Manual Go Commands / Direct Execution

-   To build:
    ```bash
    go mod tidy
    go build -o go-ssg ./cmd/main.go
    ```
-   To run (shows help):
    ```bash
    ./go-ssg
    ```
-   To generate the site:
    ```bash
    ./go-ssg generate [flags]
    ```
-   To serve the site:
    ```bash
    ./go-ssg serve [flags]
    ```
-   To edit files:
    ```bash
    ./go-ssg edit
    ```

## CLI Commands and Flags

The application now uses subcommands. Global flags:
-   `-c, --config <path>`: Path to configuration file (default: `pkg/config/config.yaml`).

### `generate`
Generates static HTML files from Markdown.
```bash
./go-ssg generate [flags]
```
**Flags for `generate`:**
-   `--contentDir <path>`: Directory containing markdown content files (overrides config).
-   `--outputDir <path>`: Directory where HTML files will be generated (overrides config).
-   `--siteTitle <title>`: Title for the site (overrides config).

### `serve`
Serves the generated static files from the output directory.
```bash
./go-ssg serve [flags]
```
**Flags for `serve`:**
-   `--outputDir <path>`: Directory of generated files to serve (overrides config, default: `output`).
-   `-p, --port <port>`: Port to serve the site on (overrides config, default: `8080`).

### `edit`
Opens a TUI to select a Markdown (from content directory) or HTML (from output directory) file to edit using an external editor (e.g., `vim`, `nano`, or `$EDITOR`).
```bash
./go-ssg edit
```
This command does not take additional flags beyond the global `--config` flag. It uses the `contentDir` and `outputDir` from the loaded configuration.

## How it Works

1.  The application loads configuration from `pkg/config/config.yaml` (or specified path).
2.  CLI flags can override these settings.
3.  The `generate` command reads Markdown (`.md`) files from the configured `contentDir`.
4.  Each Markdown file is converted into an HTML file.
5.  The resulting HTML files are saved in the configured `outputDir`.
6.  The `serve` command starts a simple HTTP server for the files from `outputDir` on the configured `serverPort`.
7.  The `edit` command provides a TUI to list and open `.md` or `.html` files in an external editor.

>## Project Structure

```
go-ssg/
├── cmd/
│   └── main.go         # Main application entry point
├── inc/
│   ├── cli/
│   │   └── cli.go      # CLI related functions (currently basic)
│   ├── config/
│   │   └── config.go   # Configuration loading and saving
│   ├── editor/
│   │   └── editor.go   # Logic for launching external editor
│   └── gen/
│       └── gen.go      # Core site generation logic
├── pkg/
│   ├── config/
│   │   └── config.yaml # Application configuration file (placeholder)
│   └── markdown/
│       └── example.md  # Example markdown (should be moved to content/)
├── content/            # (You need to create this) Directory for Markdown source files
├── output/             # Directory for generated HTML files
├── go.mod
├── go.sum
├── LICENSE
├── Makefile
└── README.md
```
