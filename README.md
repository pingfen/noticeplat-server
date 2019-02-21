# noticeplat server

## User

### create user

```bash
curl -XPOST 127.0.0.1:8080/msgpack/api/v1/users -d '
{
	"id": "b",
	"appOpenId": "b-app-openid",
	"srvOpenId": "b-srv-openid",
	"name": "B",
	"province": "Shandong",
	"phone": "11111",
	"email": "b@bingbaba.com"
}
'
```

### service binding
```bash
curl -XPOST 127.0.0.1:8080/msgpack/api/v1/users/b/service/b-srv-openid2
```

### delete user
```bash
curl -XDELETE 127.0.0.1:8080/msgpack/api/v1/users/b
```

## Message

### post message
```bash
curl -XPOST 127.0.0.1:8080/msgpack/api/v1/users/a/messages -d '
{
    "id":"000002",
    "type": 0,
    "level":3,
    "project":"MONITOR",
    "title":"机器CPU利用率过高",
    "content":"机器 127.0.0.1 CPU利用率过高,达到95%",
    "target":"a"
}
'
```

### post message comment
```bash
curl -XPOST 127.0.0.1:8080/msgpack/api/v1/users/a/messages/000002/comments -d '
{
    "content": "这是一个评论...",
    "operator":"a"
}
'
```

### list message comments
```bash
curl 127.0.0.1:8080/msgpack/api/v1/users/a/messages/000002/comments
```
