# eletrize

```
$ eletrize eletrize.json
```

```
{
  "schema": [
    {
      "label": "SCHEMA NAME",
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
        "build": {
          "method": "go",
          "args": ["build"]
        },
        "run": [
          {
            "label": "WORKER",
            "method": "./worker",
            "args": [],
            "envs": []
          }
        ]
      }
    }
  ]
}
