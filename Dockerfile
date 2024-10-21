FROM golang:1.23.2-alpine as builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o weather-bot .
RUN ls -l /app

FROM alpine:latest
WORKDIR /app
RUN apk add --no-cache ca-certificates
COPY --from=builder /app/weather-bot .
COPY .env .
# COPY wait-for-it.sh .
# RUN chmod +x wait-for-it.sh
RUN chmod +x weather-bot
CMD ["./weather-bot"]/
# CMD ["./wait-for-it.sh", "db:5432", "--", "./weather-bot"] 