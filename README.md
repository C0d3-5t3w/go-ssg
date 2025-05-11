# Go SSG (Static Site Generator)

A simple static site generator built with Go. This project takes Markdown files from a `content` directory, converts them to HTML, and serves them from an `output` directory.

## Prerequisites

- Go (version 1.24.2 or later recommended)
- Make (optional, for using the Makefile)

## Setup

1.  Clone the repository:
    ```bash
    git clone <your-repo-url>
    cd go-ssg
    ```
2.  Create a `content` directory and add your Markdown files:
    ```bash
    mkdir content
    echo "# Hello World" > content/index.md
    echo "## Another Page" > content/another.md
    ```

## Building and Running

### Using Makefile

-   To build the application:
    ```bash
    make build
    ```
-   To run the application (builds, generates site, and starts server):
    ```bash
    make run
    ```
    The site will be served at `http://localhost:8080`.

-   To clean build artifacts and the `output` directory:
    ```bash
    make clean
    ```

### Manual Go Commands

-   To build:
    ```bash
    go build -o go-ssg ./cmd/main.go
    ```
-   To run (after building):
    ```bash
    ./go-ssg
    ```

## How it Works

1.  The program reads Markdown files from the `./content/` directory.
2.  Each Markdown file is converted into an HTML file.
3.  The resulting HTML files are saved in the `./output/` directory.
4.  A simple HTTP server serves the files from the `./output/` directory.

## Project Structure

```
go-ssg/
├── cmd/
│   └── main.go         # Main application entry point
├── inc/
│   ├── cli/
│   │   └── cli.go      # CLI related functions (currently basic)
│   ├── config/
│   │   └── config.go   # Configuration loading (placeholder)
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
