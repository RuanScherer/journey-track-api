FROM golang:1.21.4-alpine

RUN apk add bash --no-cache

WORKDIR /app

ENTRYPOINT ["tail", "-f", "/dev/null"]
