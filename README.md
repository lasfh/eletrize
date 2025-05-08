# ⚡ Eletrize

[![Go Report Card](https://goreportcard.com/badge/github.com/lasfh/eletrize)](https://goreportcard.com/report/github.com/lasfh/eletrize)

**Eletrize** is a live reload tool for Go and other languages. It watches for file changes in your project and automatically runs commands, speeding up development and testing workflows.

#### 📖 English | [Português](README_ptBR.md)

---

## 🚀 Installation

Requirements:

* Go 1.23 or later

To install Eletrize, run:

```bash
go install github.com/lasfh/eletrize@latest
```

---

## ⚙️ Basic Usage

Run a simple command with file watching:

```bash
eletrize run ./server "go build" --ext=.go --label="API" --env=.env
```

This command:

* Watches the directory for changes in `.go` files.
* Runs `go build` and `.server` when changes are detected.
* Loads environment variables from the `.env` file.

---

## 📁 Configuration Files

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

## 🧪 Running a Specific Schema

If your config file contains multiple schemas, you can specify one using:

```bash
eletrize eletrize.yml --schema=1
```

Replace `1` with the desired schema index.

---

## 📝 Configuration Example

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

## 🧰 Additional Features

* **Multi-directory watching**: Define multiple schemas to monitor different parts of your project.
* **Language-agnostic support**: While optimized for Go, Eletrize can be configured for other languages.
* **Advanced customization**: Combine extensions, commands, and environment variables to tailor Eletrize to your project.

---

## 📄 License

This project is licensed under the [MIT License](LICENSE).

---
