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
FROM bash
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /app


COPY --from=build-stage app/bin/web-server ./bin/web-server
COPY --from=build-stage app/bin/migrator ./bin/migrator
COPY --from=build-stage app/scripts/ ./scripts
COPY --from=build-stage app/migrations ./migrations
RUN chmod +x ./scripts/* && chmod +x ./bin/*

USER appuser

ENTRYPOINT ["./scripts/docker-entrypoint.sh"]

