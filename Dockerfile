FROM golang:latest AS builder

WORKDIR /app/

COPY go.mod go.sum ./
RUN go mod download

RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

COPY ./ ./

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o silkroad ./cmd/app/main.go

FROM alpine:latest

COPY --from=builder /app/silkroad /app/silkroad
COPY --from=builder /app/configs /app/configs
COPY --from=builder /app/.env /app/.env
COPY --from=builder /app/db /app/db
COPY --from=builder /go/bin/migrate /usr/local/bin/migrate

WORKDIR /app

CMD ["./silkroad"]
