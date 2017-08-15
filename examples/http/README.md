# HTTP worker

## 1. Start the HTTP worker

```bash
$ go build
$ ./http
```

## 2. Send tasks via HTTP

```bash
$ curl -i -XPOST http://localhost:8080 -d '{"id": "id-1", "name": "add", "args": {"x": 1, "y": 11}}'
$ curl -i -XPOST http://localhost:8080 -d '{"id": "id-2", "name": "greet", "args": {"words": "Russell"}}'
$ curl -i -XPOST http://localhost:8080 -d '{"id": "id-3", "name": "panic", "args": {}}'
```
