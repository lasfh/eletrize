# eletrize

```
$ eletrize eletrize.json
```

```
{
  "schema": [
    {
      "name": "SCHEMA NAME",
      "envs": {
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
      "commands": {
        "run": [
          {
            "name": "WORKER",
            "method": "go",
            "args": ["run", "main.go"],
            "envs": []
          }
        ]
      }
    }
  ]
}
