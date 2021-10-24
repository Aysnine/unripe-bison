# unripe-bison

- [GoFiber](https://docs.gofiber.io/)
- [CockroachDB (Serverless)](https://www.cockroachlabs.com/)

``` bash
# Running at http://127.0.0.1:9000
DATABASE_CONNECTION='<database connect string>' go run app.go
```

API Routes:

- `GET /` Root check heath
- `GET /dashboard` Fiber dashboard
- `GET /api/books` Books example
