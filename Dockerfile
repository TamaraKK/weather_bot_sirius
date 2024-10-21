FROM golang:1.23.2 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o myapp .

FROM alpine:latest

RUN apk add --no-cache ca-certificates

COPY --from=builder /app/myapp .

CMD ["./myapp"]