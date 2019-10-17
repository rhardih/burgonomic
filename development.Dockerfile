FROM golang:1.12-alpine

RUN apk update && apk add --no-cache git

WORKDIR /app

COPY . .

RUN go build -o main .

EXPOSE 8080
