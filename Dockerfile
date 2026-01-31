FROM golang:1.24 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o t-log .


FROM debian:bookworm-slim

WORKDIR /app

COPY --from=builder /app/t-log .

RUN mkdir -p /app/output


ENTRYPOINT ["./t-log"]
