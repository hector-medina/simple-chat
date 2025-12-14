ARG GO_VERSION=1.25.5

FROM golang:${GO_VERSION}-bookworm AS builder
RUN apt-get update \
    && apt-get install -y --no-install-recommends pkg-config libzmq3-dev \
    && rm -rf /var/lib/apt/lists/*
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=1 GOOS=linux go build -o /out/server ./server.go
RUN CGO_ENABLED=1 GOOS=linux go build -o /out/client ./client.go

FROM debian:bookworm-slim AS server
RUN apt-get update \
    && apt-get install -y --no-install-recommends libzmq5 \
    && rm -rf /var/lib/apt/lists/*
COPY --from=builder /out/server /usr/local/bin/server
EXPOSE 5555 5563
ENTRYPOINT ["/usr/local/bin/server"]

FROM debian:bookworm-slim AS client
RUN apt-get update \
    && apt-get install -y --no-install-recommends libzmq5 \
    && rm -rf /var/lib/apt/lists/*
COPY --from=builder /out/client /usr/local/bin/client
ENTRYPOINT ["/usr/local/bin/client"]
