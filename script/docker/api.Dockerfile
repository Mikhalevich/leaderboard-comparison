FROM golang:1.25-alpine3.22 AS builder

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -mod=vendor -a -installsuffix cgo -ldflags="-w -s" -o ./bin/api cmd/api/main.go

FROM alpine:3.22

EXPOSE 8080

WORKDIR /app/

COPY --from=builder /app/bin/api /app/api
COPY --from=builder /app/config/config-api.yaml /app/config-api.yaml

ENTRYPOINT ["./api", "-config", "config-api.yaml"]
