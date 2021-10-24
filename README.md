# unripe-bison

- [GoFiber](https://docs.gofiber.io/)
- [CockroachDB (Serverless)](https://www.cockroachlabs.com/)

## Quick start

``` bash
# Running at http://127.0.0.1:9000
DATABASE_CONNECTION='<database connect string>' go run app.go
```

Use `.env` file:

``` bash
DATABASE_CONNECTION=postgresql://xxx
```

## API Routes

- `GET /` Root check heath
- `GET /dashboard` Fiber dashboard
- `GET /api/books` Books example
