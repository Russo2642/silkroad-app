FROM golang:latest

ENV GOPATH=/

WORKDIR /app/

COPY ./ ./

# install psql
RUN apt-get update
RUN apt-get -y install postgresql-client

# make wait-for-postgres.sh executable
RUN chmod +x wait-for-postgres.sh

# build go app
RUN go mod download
RUN go build -o silkroad ./cmd/app/main.go

CMD ["./silkroad"]
