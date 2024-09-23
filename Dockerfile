FROM golang:latest AS builder

WORKDIR /app/

COPY go.mod go.sum ./
RUN go mod download

COPY ./ ./

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o silkroad ./cmd/app/main.go

# Финальный образ
FROM alpine:latest

RUN apk add --no-cache postgresql-client

COPY --from=builder /app/silkroad /app/silkroad
COPY --from=builder /app/configs /app/configs
COPY --from=builder /app/.env /app/.env

COPY --from=builder /app/wait-for-postgres.sh /app/wait-for-postgres.sh
RUN chmod +x /app/wait-for-postgres.sh

WORKDIR /app

CMD ["./silkroad"]
