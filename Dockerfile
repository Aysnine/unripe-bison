FROM golang:alpine AS build

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build -o main .

FROM scratch

WORKDIR /app

COPY --from=build /app/main /app/main
COPY --from=build /app/public /app/public

EXPOSE 9000

ENTRYPOINT ["/app/main"]
