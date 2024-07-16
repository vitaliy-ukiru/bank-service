#FROM ubuntu:latest
#LABEL authors="Vitaliy"
#
#ENTRYPOINT ["top", "-b"]

# syntax=docker/dockerfile:1

# Build the application from source
FROM golang:1.22 AS build-stage

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux go build -o ./bin/migrator github.com/vitaliy-ukiru/bank-service/cmd/migrator
RUN CGO_ENABLED=0 GOOS=linux go build -o ./bin/web-server github.com/vitaliy-ukiru/bank-service/cmd/api



# Deploy the application binary into a lean image
FROM scratch
RUN addgroup --system app && adduser --system --group app

WORKDIR /app

COPY --from=build-stage app/bin/web-server ./bin/web-server
COPY --from=build-stage app/bin/migrator ./bin/migrator
RUN chmod +x scripts/* && chmod +x ./bin/*

USER app

ENTRYPOINT ["./scripts/docker-entrypoint.sh"]

