# ⚡ Eletrize

[![Go Report Card](https://goreportcard.com/badge/github.com/lasfh/eletrize)](https://goreportcard.com/report/github.com/lasfh/eletrize)

**Eletrize** é uma ferramenta de *live reload* para projetos em Golang e outras linguagens. Ela monitora alterações nos arquivos do seu projeto e executa comandos automaticamente, agilizando o desenvolvimento e os testes.

### 📖 [English](./README.md) | Português  

---

## 🚀 Instalação

Requisitos:

* Go 1.23 ou superior

Para instalar o Eletrize, execute:

```bash
go install github.com/lasfh/eletrize@latest
```

---

## ⚙️ Uso Básico

Execute um comando simples com monitoramento de alterações:

```bash
eletrize run ./server "go build" --ext=.go --label="API" --env=.env
```

Este comando:

* Monitora o diretório por alterações em arquivos `.go`.
* Executa `go build` e `./server` sempre que uma alteração é detectada.
* Utiliza as variáveis de ambiente definidas no arquivo `.env`.

---

## 📁 Arquivos de Configuração

O Eletrize pode detectar automaticamente arquivos de configuração com os seguintes nomes:

* `eletrize.yml`
* `eletrize.yaml`
* `.eletrize.yml`
* `.eletrize.yaml`
* `eletrize.json`
* `.eletrize.json`
* `.eletrize` (formato JSON)

Para executar o Eletrize com um arquivo de configuração:

```bash
eletrize
```

Ou especifique o caminho do arquivo:

```bash
eletrize path/eletrize.yml
```

---

## 🧪 Executando com um Schema Específico

Se o seu arquivo de configuração contém múltiplos schemas, você pode especificar qual deseja executar:

```bash
eletrize eletrize.yml --schema=1
```

Substitua `1` pelo número correspondente ao schema desejado.

---

## 📝 Exemplo de Arquivo de Configuração

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

## 🔧 Configuração do VSCode Launch

O Eletrize pode detectar e utilizar automaticamente as configurações de launch do VSCode a partir do arquivo `.vscode/launch.json`. Esta funcionalidade permite aproveitar suas configurações de debug existentes no VSCode para live reloading.

Para usar a detecção automática do VSCode launch:

```bash
eletrize
```

O Eletrize detectará automaticamente:

* Configurações de launch do Go com `"type": "go"`, `"request": "launch"` e `"mode": "auto"`
* Caminho do programa (suporta a variável `${workspaceFolder}`)
* Variáveis de ambiente e arquivos de ambiente
* Argumentos de linha de comando

**Exemplo de `.vscode/launch.json`:**

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

Esta configuração será automaticamente convertida para monitorar arquivos `.go` e executar o servidor com live reloading.

---

## 🧰 Recursos Adicionais

* **Monitoramento de múltiplos diretórios**: Configure vários schemas para monitorar diferentes partes do seu projeto simultaneamente.
* **Suporte a múltiplas linguagens**: Embora otimizado para Golang, o Eletrize pode ser configurado para outras linguagens.
* **Personalização avançada**: Combine diferentes extensões, comandos e variáveis de ambiente para adaptar o Eletrize às necessidades específicas do seu projeto.
* **Integração com VSCode**: Detecta e utiliza automaticamente configurações de launch do VSCode para um fluxo de desenvolvimento seamless.

---

## 📄 Licença

Este projeto está licenciado sob a [Licença MIT](LICENSE).

---
