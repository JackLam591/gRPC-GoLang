FROM golang:1.18.3-alpine AS builder

RUN mkdir /app
WORKDIR /app

RUN apk add git libc-dev build-base
COPY go.mod go.sum ./
RUN go mod download

COPY . ./
RUN go build -o /app/go-gRPC

FROM alpine:3.14
RUN apk update \
    && apk add ca-certificates \
    && rm -rf /var/cache/apk/*
WORKDIR /app

COPY --from=builder /app/go-gRPC /app/go-gRPC
# COPY .env ./

# CMD printenv

CMD ["./go-gRPC"]