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
The servers are configured through environment variables. The following variables are available

Name | Default | Description
----|------|-------
PALLAS_GRPC_SERVE_ENDPOINT | :50051 |
PALLAS_DB_SERVER | db |
PALLAS_DB_PORT | 5432 |
PALLAS_DB_NAME | gotp |
PALLAS_DB_USER | postgres |
PALLAS_DB_PASS | passpass |
PALLAS_GRPC_ENDPOINT | server:50051" |
PALLAS_REST_SERVE_ENDPOINT | :8081 |

## Requirements
Docker >= 19.03.0

## Development
Golang >= 1.18.0
