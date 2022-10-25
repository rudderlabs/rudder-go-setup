# syntax=docker/dockerfile:1
ARG GO_VERSION={{.GoVersion}}
ARG ALPINE_VERSION=3.16
FROM golang:${GO_VERSION}-alpine${ALPINE_VERSION} AS builder
ARG VERSION
ARG REVISION


COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN mkdir /builds

RUN for d in */; do; go build -o /builds/${d%/} ./${d}; done

FROM alpine:${ALPINE_VERSION}

RUN apk -U --no-cache upgrade && \
    apk add --no-cache tzdata ca-certificates postgresql-client curl bash


COPY --from=builder /builds .