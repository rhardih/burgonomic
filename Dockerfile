FROM golang:1.12-alpine AS build

RUN apk update && apk add --no-cache git

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 GOOS=linux \
  go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o main .

FROM scratch

COPY --from=build /app/main ./main
COPY --from=build /app/html ./html
COPY --from=build /app/static ./static
COPY --from=build /app/big-mac-adjusted-index.csv ./big-mac-adjusted-index.csv

EXPOSE 8080

CMD ["./main"]
