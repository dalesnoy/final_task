FROM golang:1.21 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /todo-app main.go
FROM ubuntu:latest
WORKDIR /app
COPY --from=builder /todo-app .
COPY web ./web

ENV TODO_PORT=7540
ENV TODO_DBFILE=/app/data/scheduler.db
EXPOSE 7540
CMD ["./todo-app"]
