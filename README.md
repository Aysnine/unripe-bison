# unripe-bison

**👏 Available Online DEMO: https://unripe-bison.cnine.me**
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

## API Routes

- `GET /` Root check heath
- `GET /dashboard` Fiber dashboard
- `GET /api/books` Books example
