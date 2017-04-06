# NSQ worker

## 1. Start NSQD

```bash
$ nsqd
```

## 2. Start the NSQ worker

```bash
$ go build
$ ./nsq --nsqd-tcp-address="127.0.0.1:4150" --topic="test"
```

## 3. Send tasks via NSQ

```bash
$ curl http://127.0.0.1:4151/pub?topic=test -d '{"uuid": "uuid-1", "name": "add", "args": {"x": 1, "y": 11}}'
$ curl http://127.0.0.1:4151/pub?topic=test -d '{"uuid": "uuid-2", "name": "greet", "args": {"words": "Russell"}}'
$ curl http://127.0.0.1:4151/pub?topic=test -d '{"uuid": "uuid-3", "name": "panic", "args": {}}'
```
