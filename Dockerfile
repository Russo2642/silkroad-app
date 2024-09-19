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



# Stage 1: Build the Go application
#FROM golang:latest AS builder
#
#WORKDIR /app/
#
## Copy the Go modules and download dependencies
#COPY go.mod go.sum ./
#RUN go mod download
#
## Copy the rest of the application code
#COPY ./ ./
#
## Build the Go app
#RUN go build -o silkroad ./cmd/app/main.go
#
## Stage 2: Prepare the final minimal image
#FROM alpine:latest
#
## Install PostgreSQL client
#RUN apk update && apk add --no-cache postgresql-client
#
#WORKDIR /app/
#
## Copy the built Go binary from the builder stage
#COPY --from=builder /app/silkroad .
#
## Copy wait-for-postgres.sh script and make it executable
##COPY --from=builder /app/wait-for-postgres.sh /app/wait-for-postgres.sh
##RUN chmod +x /app/wait-for-postgres.sh
#
## Command to run the Go app
#CMD ["./silkroad"]

