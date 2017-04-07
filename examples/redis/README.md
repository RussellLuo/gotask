# Redis worker

## 1. Start redis-server

```bash
$ redis-server
```

## 2. Start the Redis worker

```bash
$ go build
$ ./redis -queue="test"
```

## 3. Send tasks via redis-cli

```bash
$ redis-cli
127.0.0.1:6379> RPUSH test '{"uuid": "uuid-1", "name": "add", "args": {"x": 1, "y": 11}}'
127.0.0.1:6379> RPUSH test '{"uuid": "uuid-2", "name": "greet", "args": {"words": "Russell"}}'
127.0.0.1:6379> RPUSH test '{"uuid": "uuid-3", "name": "panic", "args": {}}'
```