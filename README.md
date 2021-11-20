# unripe-bison

*Web server example by GoFiber⚡️ and CockroachDB📖. Just a random project name😜.*

**🎉🎉🎉 Online DEMO (Kubernetes): https://unripe-bison.cnine.me**

``` bash
# You can try this
curl https://unripe-bison.cnine.me/api/books

# Performance testing eg 233 times
curl https://unripe-bison.cnine.me/api/books\?_times\=\[0-233\] -w "%{http_code} total:%{time_total}s size:%{size_download}\n" -o /dev/null -s

# Websocket global chat by https://github.com/vi/websocat
websocat "wss://unripe-bison.cnine.me/chat?nickname=Jerry"
```

Docs References

- [GoFiber](https://docs.gofiber.io/)
- [CockroachDB (Serverless)](https://www.cockroachlabs.com/)

**Page Routes**

- `/` Home check heath
- `/chat` WebSocket chat room
- `/monitor` Fiber monitor
- `/swagger/doc.json` Swagger document data only

**API Routes**

- `GET /api/books` Books example
- `GET /api/hongkong-weather` Request from this [Free API](https://data.weather.gov.hk/weatherAPI/opendata/weather.php?dataType=fnd&lang=sc)
- `GET /api/random-anime-image` Request from this [WaiFu.PICS](https://waifu.pics/docs)

## Quick start

Add `.env` file and set variables:

``` bash
MODE=development
DATABASE_CONNECTION=postgresql://your-database-connection-string-here
CHAT_REDIS_CONNECTION=redis://your-redis-connection-string-here
```

- MODE = `local` | `development` | `production`

Command:

``` bash
# Running at http://127.0.0.1:9000
go run app.go

# Quick Unit test
go test ./...

# Make swagger docs when API update.
# Need global install: https://github.com/arsmn/fiber-swagger
swag init
```

NodeJS Helper:

``` bash
# Before each command
yarn install

# Make a release commit
yarn release
```

## TODO

- Books full CRUD
- Websocket Demo
- Custom Error response
- Login auth
- Transactions Demo
- Unit test with mocked database. eg: [pgmock](https://github.com/jackc/pgmock)
- Auto run initial sql
- Auto Migrations

## Learn References

- [A Complete Guide to JSON in Golang (With Examples)](https://www.sohamkamani.com/golang/json/)
- [Go 语言中的面向对象](http://kangkona.github.io/oo-in-golang/)
- [Go 语言高性能编程](https://github.com/geektutu/high-performance-go)
- [Error handling and go](https://go.dev/blog/error-handling-and-go)
- [Scaling Websockets in the Cloud (Part 1). From Socket.io and Redis to a distributed architecture with Docker and Kubernetes](https://dev.to/sw360cab/scaling-websockets-in-the-cloud-part-1-from-socket-io-and-redis-to-a-distributed-architecture-with-docker-and-kubernetes-17n3)
- [Scaling websockets](https://github.com/sw360cab/websockets-scaling)
- [Building Chat Service in Golang and Websockets Backed by Redis](https://levelup.gitconnected.com/building-chat-service-in-golang-and-websockets-backed-by-redis-b42a8784636c)
