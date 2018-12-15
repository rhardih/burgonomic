FROM golang:alpine AS build

RUN apk add --update \
  git

WORKDIR /go/src/burger-pricing

COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

FROM alpine:3.7

RUN apk add --no-cache ca-certificates
COPY --from=build /go/src/burger-pricing/html ./html
COPY --from=build /go/bin/burger-pricing .

EXPOSE 8080

CMD ["./burger-pricing"]
