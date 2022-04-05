# pallas
**This is a Babygopher project in its early alpha stage!**
**Please do not run this as production service!**

Self-hosted TOTP and HOTP sync suite based on Go.
This suite defines the protocol through Protocol Buffers and exposes gRPC and REST servers.

### Planned Features
- [ ] Improved logging
- [ ] Secret column encryption in DB
- [ ] Authentication
- [ ] Multi tenant support
- [ ] Multi vault support

## Quickstart
To get up and running, run

```docker-compose up -d```

in the project root.

The servers will be available on ports 50051 (gRPC) and 8001 (REST).

**NOTE:** The service currently neither features HTTPS nor authentication as it's in POC phase!

## Configuration
The application is configured through environment variables. These are primarily configured through the `.env` file in the project root.
If you intend to run the binaries on bare metal, remember to set these variables in your OS accordingly.

Name | Default | Component | Description
----|------|-------|-----
PALLAS_GRPC_SERVE_ENDPOINT | :50051 | server | Endpoint the gRPC API is served at
PALLAS_DB_SERVER | db | server | DB host the server connects to
PALLAS_DB_PORT | 5432 | server | DB port the server connects to
PALLAS_DB_NAME | pallas | server | DB name the server uses
PALLAS_DB_USER | postgres | server & db | DB username
PALLAS_DB_PASS | passpass | server & db | DB password for given user
PALLAS_GRPC_ENDPOINT | server:50051 | rest-api | Connection string the rest-api proxies
PALLAS_REST_SERVE_ENDPOINT | :8081 | rest-api | Endpoint the REST API is served at

## Requirements
Docker >= 19.03.0

## Development
Golang >= 1.18.0
