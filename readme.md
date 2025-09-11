# message-responder

Service that consumes normalized Telegram updates from Kafka and responds with guidance or forwards image OCR requests.

## Requirements

- Go 1.23+

## Env

- MSG_RESP_KAFKA_BOOTSTRAP_SERVERS_VALUE — Kafka broker, e.g. `localhost:9092`
- MSG_RESP_KAFKA_GROUP_ID — consumer group id
- MSG_RESP_KAFKA_REQUEST_TOPIC_NAME — incoming requests topic
- MSG_RESP_KAFKA_RESPONSE_TOPIC_NAME — outgoing responses topic
- MSG_RESP_KAFKA_OCR_TOPIC_NAME — OCR requests topic
- MSG_RESP_KAFKA_SASL_USERNAME / MSG_RESP_KAFKA_SASL_PASSWORD — optional SASL/SCRAM creds
- MSG_RESP_KAFKA_CLIENT_ID — client id

## Run

- set -a && source .env && set +a && go run ./cmd/message-responder
- or: export $(cat .env | xargs) && go run ./cmd/message-responder

## Setup

- go mod tidy

## Build

- go build -v -o bin/message-responder ./cmd/message-responder

## Docker

- docker buildx build --no-cache --progress=plain .

## Release (tags)

- git tag v0.1.1 && git push origin v0.1.1

## Notes

- SASL is enabled only if username or password is set.
- Logging uses Zap production logger.

## Command Guide

### Run with exported .env (one‑liner)

Exports all variables from `.env` into the current shell and runs the service.

```
export $(cat .env | xargs) && go run ./cmd/message-responder
```

### Run with `source` (safer for complex values)

Loads `.env` preserving quotes and special characters, then runs the service.

```
set -a && source .env && set +a && go run ./cmd/message-responder
```

### Fetch/clean module deps

Resolves dependencies and prunes unused ones.

```
go mod tidy
```

### Verbose build (diagnostics)

Builds the binary with verbose and command tracing. Removes old binary after build to keep the tree clean.

```
go build -v -x ./cmd/message-responder && rm -f message-responder
```

### Docker build (Buildx)

Builds the image with detailed progress logs and without cache.

```
docker buildx build --no-cache --progress=plain .
```

### Create and push tag

Cuts a release tag and pushes it to remote.

```
git tag v0.0.1
git push origin v0.0.1
```

### Manage tags

List all tags, delete a tag locally and remotely, verify deletion.

```
git tag -l
git tag -d vX.Y.Z
git push --delete origin vX.Y.Z
git ls-remote --tags origin | grep 'refs/tags/vX.Y.Z$'
```
