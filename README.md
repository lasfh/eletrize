# eletrize

```
$ eletrize eletrize.json
```

```
{
  "schema": [
    {
      "name": "SCHEMA NAME",
      "env": {
        "key": "value"
      },
      "watcher": {
        "path": ".",
        "recursive": true,
        "extensions": [
          ".go",
          ".json"
        ]
      },
      "command": {
        "name": "go",
        "args": ["run", "main.go"]
      }
    }
  ]
}
