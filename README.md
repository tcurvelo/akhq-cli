# akhq-cli

A command-line client for [AKHQ](https://akhq.io).

## Use cases

- **Watching offset lag**: poll a topic's offset lag per consumer group at a finer interval than metrics scraping allows (see `sampler.yml`).
- **Topic dumps**: fetch every message of a topic, e.g. a dead letter topic, for a report or a manual replay (planned).
- **Finding specific messages**: search a topic by key, value or header, optionally starting from a timestamp.

## Install

```sh
go install github.com/tcurvelo/akhq-cli@latest
```

Requires Go 1.27+.

## Configuration

Set via environment variables: `AKHQ_URL`, `AKHQ_CLUSTER`, and either `AKHQ_TOKEN` or `AKHQ_USER`/`AKHQ_PASS`. Set `AKHQ_INSECURE=true` to skip TLS verification.

## Usage

```sh
akhq-cli groups <topic>   # total offset lag per consumer group
akhq-cli search <topic> --search-by-key <key>
```

## Development

```sh
make build test lint
```

## TODO

- [ ] [sampler] Authenticate first, then reuse the token
- [ ] [akhq] Add command for authenticating only
- [ ] [sampler] Get all consumer groups in a same request
