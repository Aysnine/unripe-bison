# unripe-bison

*Web server example by GoFiber⚡️ and CockroachDB📖. Just a random project name😜.*

**🎉🎉🎉 Online DEMO: https://unripe-bison.cnine.me**

``` bash
# You can try this
curl https://unripe-bison.cnine.me/api/books
```

Docs References

- [GoFiber](https://docs.gofiber.io/)
- [CockroachDB (Serverless)](https://www.cockroachlabs.com/)

## Quick start

Add `.env` file and set variables:

``` bash
DATABASE_CONNECTION=postgresql://your-database-connection-string-here
```

Command:

``` bash
# Running at http://127.0.0.1:9000
go run app.go

# Quick Unit test
go test ./...
```

NodeJS Helper:

``` bash
# Before each command
yarn install

# Make a release commit
yarn release
```

## API Routes

- `GET /` Root check heath
- `GET /dashboard` Fiber dashboard
- `GET /hongkong-weather` Request from this [Free API](https://data.weather.gov.hk/weatherAPI/opendata/weather.php?dataType=fnd&lang=sc)
- `GET /api/books` Books example

## TODO

- Books full CRUD
- Websocket Demo
- Custom Error response
- Login auth
- Transactions Demo
- Unit test with mocked database. eg: [pgmock](https://github.com/jackc/pgmock)
- Auto run initial sql
- Auto Migrations
