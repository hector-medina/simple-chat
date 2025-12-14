# Chat Project Reproduction Prompt


```
Create a Go project named “chat” with the precise structure and contents below. Use Go 1.25.5, module path `chat`, and only depend on `github.com/pebbe/zmq4 v1.2.9` (include the matching `go.sum` entries).

Project layout:

.
├── client/
│   ├── client_comm.go
│   └── participant.go
├── server/
│   ├── chat.go
│   ├── message_store.go
│   └── server.go
├── shared/message.go
├── client.go
├── server.go
├── README.md
├── Dockerfile
├── docker-compose.yml
└── .gitignore  (contains the single line `*.DS_Store`)

### shared/message.go
Define struct `Message` with fields:
```
ID      int    `json:"id"`
Channel string `json:"channel"`
Author  string `json:"author"`
Message string `json:"message"`
```
Keep the existing logical-design comments describing the data container.

### client/participant.go
Implement type `Participant` with private fields `nick string`, `channel string`, `mu sync.Mutex`. Methods:

- `NewParticipant(nick, channel string) *Participant`
- `Channel() string`
- `TextRead(text string)` creates a `shared.Message` and invokes `SendMessage`.
- `MessageArrived(msg shared.Message)` prints `[channel][author]: message` while holding the mutex.

Preserve the logical-design comments describing constructor, `text_read`, and `message_arrived`.

### client/client_comm.go
Imports: `encoding/json`, `fmt`, `log`, `os`, `sync`, `chat/shared`, `github.com/pebbe/zmq4`.

Define:
```
const defaultServerHost = "localhost"

var (
    serverHost        = getServerHost()
    pushServerAddress = fmt.Sprintf("tcp://%s:5555", serverHost)
    pubServerAddress  = fmt.Sprintf("tcp://%s:5563", serverHost)
)
```

`getServerHost()` returns `os.Getenv("CHAT_SERVER_HOST")` or `defaultServerHost` if unset.

Maintain a package-level PUSH socket singleton using `sync.Once`. `SendMessage` marshals the `shared.Message` to JSON and sends it via a `zmq.PUSH` connected to `pushServerAddress`, logging any errors. `CheckMessages` creates a `zmq.SUB`, connects to `pubServerAddress`, subscribes to `p.Channel()`, then in an infinite loop receives the envelope (ignored) and payload frame, unmarshals into `shared.Message`, and calls `p.MessageArrived`. Log all errors with `log.Printf`.

### client.go
Main program: default `name := "Anonymous"`, `channel := "general"`, override with `os.Args[1]`/`os.Args[2]` if provided. Print participant name/channel, create the participant, start `go client.CheckMessages(p)`, and read stdin with `bufio.Scanner`, calling `p.TextRead` for each line.

### server/message_store.go and server/chat.go
`MessageStore` holds a slice of messages, an incrementing `nextID`, and a mutex. `Add(msg shared.Message) shared.Message` assigns the next ID, appends, and returns the stored message. `FetchAfter(lastID int, channel string)` returns all messages from that channel with `ID > lastID`. `Chat` wraps the store; `DistributeMessage` logs `channel/author/msg` and calls `store.Add`, returning the stored message.

### server/server.go
Imports `encoding/json`, `log`, `chat/shared`, `github.com/pebbe/zmq4`. Constants:
```
const (
    pullBindAddress = "tcp://*:5555"
    pubBindAddress  = "tcp://*:5563"
)
```
`StartServer()` creates the store/chat, opens a `zmq.PULL` socket bound to `pullBindAddress` and a `zmq.PUB` bound to `pubBindAddress`, logging fatal errors during setup. Loop forever: `pull.RecvBytes`, `json.Unmarshal`, pass to `chat.DistributeMessage`, marshal the stored message, send the channel as the PUB envelope (`zmq.SNDMORE`) and the JSON payload. Log any processing errors.

### server.go
`main()` prints `Starting server...` and invokes `server.StartServer()`.

### README.md
Exact contents:

```
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
```

### Dockerfile
Use:

```
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
```

### docker-compose.yml
No `version` key. Content:

```
services:
  server:
    build:
      context: .
      target: server
    ports:
      - "5555:5555"
      - "5563:5563"

  client:
    build:
      context: .
      target: client
    depends_on:
      - server
    stdin_open: true
    tty: true
    environment:
      - CHAT_SERVER_HOST=server
```

This prompt must reproduce every file and behavior exactly as described.
```
