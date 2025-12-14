# simple-chat

This repository contains a simple chat implementation for educational purposes, not ready for production deployment.

## Running with Docker

The repository ships with a multi-stage `Dockerfile` and `docker-compose.yml` that wrap the ZeroMQ server and the CLI client.

```bash
# build the images
docker compose build

# run the server in the background
docker compose up server

# in a separate terminal you can attach clients (they auto-connect to the server container)
docker compose run --rm client                # Anonymous/general
docker compose run --rm client Alice general  # Custom nickname/channel
```

The client container keeps STDIN/STDOUT attached, so you can type messages interactively. Launch multiple client containers (possibly specifying different channels) to chat back and forth.

## Running locally

If you prefer to run the binaries on your host you will need the ZeroMQ headers and runtime libraries (`pkg-config` and `libzmq`). On macOS, for example:

```bash
brew install pkg-config zeromq
go run server.go            # terminal 1
go run client.go Bob games  # terminal 2
```

Use additional terminals for more participants.
