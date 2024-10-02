FROM golang:latest AS builder

WORKDIR /app/

COPY go.mod go.sum ./
RUN go mod download

COPY ./ ./

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o silkroad ./cmd/app/main.go

FROM alpine:latest

COPY --from=builder /app/silkroad /app/silkroad
COPY --from=builder /app/configs /app/configs

WORKDIR /app

CMD ["./silkroad"]