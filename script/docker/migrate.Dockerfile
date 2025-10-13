FROM golang:1.25-alpine3.22 AS builder

WORKDIR /app

RUN GOBIN=/app go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@v4.19.0

FROM alpine:3.22

WORKDIR /app/

COPY --from=builder /app/migrate /app/migrate
COPY script/db/postgres/migrations /app/db/migrations

ENTRYPOINT ["./migrate", "-database", "postgres://leaderboard:leaderboard@postgres:5432/leaderboard?sslmode=disable", "-path", "db/migrations", "up"]
