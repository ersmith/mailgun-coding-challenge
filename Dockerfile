FROM golang:1.15-alpine3.13 AS builder

COPY . /app
WORKDIR /app

ENV HTTP_PORT=8080

RUN go mod download
RUN go build -o /app/mailgun-challenge .

FROM alpine:3.13

RUN adduser -D mailgun
USER mailgun

COPY --from=builder /app /home/mailgun/app

WORKDIR /home/mailgun/app

ENV HTTP_PORT=8080

EXPOSE 8080

CMD ["./mailgun-challenge"]
