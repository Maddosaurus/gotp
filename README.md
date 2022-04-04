# gotp
**This is a Babygopher project in its early alpha stage!**
**Please do not run this as production service!**

Self-hosted TOTP and HOTP sync suite based on Go.
This suite defines the protocol through Protocol Buffers and exposes gRPC and REST servers.

## Quickstart
To get up and running, run

```docker-compose up -d```

in the project root.

The servers will be available on ports 50051 (gRPC) and 8001 (REST).

**NOTE:** The service currently neither features HTTPS nor authentication as it's in POC phase!

## Requirements
Docker >= 19.03.0

## Development
Golang >= 1.18.0
