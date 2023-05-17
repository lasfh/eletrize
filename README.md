# eletrize - Golang live reload

```
$ eletrize eletrize.json
```

```
{
  "schema": [
    {
      "label": "SCHEMA NAME",
      "ignore_notification": false,
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
          "args": ["build"],
          "envs": {}
        },
        "run": [
          {
            "label": "WORKER",
            "method": "./worker",
            "args": [],
            "envs": {}
          }
        ]
      }
    }
  ]
}
