FROM golang:latest

# Install air
RUN go install github.com/cosmtrek/air@latest

# Intall go migrate tool
RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

RUN mkdir main

WORKDIR /app