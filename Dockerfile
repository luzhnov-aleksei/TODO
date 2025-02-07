FROM golang:1.23.5-alpine AS builder

WORKDIR /app
COPY . .
RUN go build -o todo .

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/todo .

EXPOSE 3000
CMD ["./todo"]