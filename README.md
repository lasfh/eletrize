# Eletrize

[![Go Report Card](https://goreportcard.com/badge/github.com/lasfh/eletrize)](https://goreportcard.com/report/github.com/lasfh/eletrize)

**Eletrize** is a live reload tool for Go and other languages. It watches for file changes in your project and automatically runs commands, speeding up development and testing workflows.

#### English | [Português](README_ptBR.md)

---

## Installation

Requirements:

* Go 1.23 or later

To install Eletrize, run:

```bash
go install github.com/lasfh/eletrize@latest
```

---

## Basic Usage

Run a simple command with file watching:

```bash
eletrize run [run] [build] [flags]
# Example:
eletrize run ./server "go build -o server" --ext=.go,.mod --label="API" --env=.env
```

This command:

* Watches the directory for changes in `.go` and `.mod` files.
* Runs `go build -o server` and `./server` when changes are detected.
* Loads environment variables from the `.env` file.

---

## Configuration Files

Eletrize can automatically detect configuration files named:

* `eletrize.yml`
* `eletrize.yaml`
* `.eletrize.yml`
* `.eletrize.yaml`
* `eletrize.json`
* `.eletrize.json`
* `.eletrize` (JSON format)

To run Eletrize with a config file:

```bash
eletrize
```

Or specify the path manually:

```bash
eletrize path/eletrize.yml
```

---

## Running a Specific Schema

If your config file contains multiple schemas, you can specify one using:

```bash
eletrize eletrize.yml --schema=1
```

Replace `1` with the desired schema index.

---

## Configuration Example

```yaml
schema:
  - label: API
    workdir: ""
    envs:
      key: "value"
    env_file: ".env"
    watcher:
      path: "."
      recursive: true
      excluded_paths:
        - "frontend"
      extensions:
        - ".go"
    commands:
      build:
        method: "go"
        args:
          - "build"
        envs:
          key: "value"
        env_file: ""
      run:
        - method: "./server"
          envs:
            PORT: "8080"
          env_file: ""
```

---

## VSCode Launch Configuration

Eletrize can automatically detect and use VSCode launch configurations from `.vscode/launch.json`. This feature allows you to leverage your existing VSCode debug configurations for live reloading.

To use VSCode launch detection:

```bash
eletrize
```

Eletrize will automatically detect:

* Go launch configurations with `"type": "go"`, `"request": "launch"`, and `"mode": "auto"`
* Program path (supports `${workspaceFolder}` variable)
* Environment variables and environment files
* Command line arguments

**Example `.vscode/launch.json`:**

```json
{
    "configurations": [
        {
            "name": "Launch Server",
            "type": "go",
            "request": "launch",
            "mode": "auto",
            "program": "${workspaceFolder}/cmd/server",
            "args": ["--port", "8080"],
            "envFile": "${workspaceFolder}/.env",
            "env": {
                "DEBUG": "true"
            }
        }
    ]
}
```

This configuration will be automatically converted to watch `.go` files and run the server with live reloading.

---

### `help`
Display help information.
```bash
eletrize help [command]
```

---

## Comparison: Eletrize vs Air

Both tools are great for live reloading, but they have different focuses:

| Feature | Eletrize | Air |
|---------|----------|-------|
| **Language Support** | **Agnostic** (Go, Rust, Node, etc) | Go Focused |
| **VSCode Integration** | **Native** (Reads `launch.json`) | Manual config required |
| **Configuration** | YAML, JSON (Multiple schemas) | TOML |
| **Multi-folder** | **Yes** (Monorepo readiness) | Limited |

**Why choose Eletrize?**
If you work with multiple languages or want zero-config integration with your VSCode debugger, Eletrize is the way to go.

---

## License

This project is licensed under the [MIT License](LICENSE).
