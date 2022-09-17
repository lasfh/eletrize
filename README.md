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
        "build": {
          "method": "go",
          "args": ["build"]
        },
        "run": [
          {
            "name": "WORKER",
            "method": "./worker",
            "args": [],
            "envs": []
          }
        ]
      }
    }
  ]
}
