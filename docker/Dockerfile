FROM golang:alpine AS builder

LABEL authors="Ahmad Rifa'i"

WORKDIR /app

RUN apk add --no-cache make bash

COPY ../go.mod go.sum ./
RUN go mod download

COPY .. .

RUN --mount=type=cache,target=/root/.cache \
    --mount=type=cache,target=/go/pkg/mod \
    make build

FROM alpine:latest

WORKDIR /app

RUN apk add --no-cache tzdata ca-certificates \
    && addgroup -S app \
    && adduser -S setetes -G app

ENV TZ=Asia/Jakarta
ENV GIN_MODE=release

COPY --from=builder /bin/app ./app

USER setetes

EXPOSE 8080

ENTRYPOINT ["./app"]
