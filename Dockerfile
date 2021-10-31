FROM golang:alpine AS build

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build -o /main .

FROM alpine

WORKDIR /

COPY --from=build /main /main
COPY --from=build /app/public /public

RUN adduser -S -D -H -h /app appuser
USER appuser

EXPOSE 9000

CMD ["./main"]
