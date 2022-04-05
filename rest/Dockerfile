# syntax=docker/dockerfile:1

#
# Build container
#
FROM golang:1.18-alpine as builder

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /build

COPY ./go.mod ./go.sum ./
RUN go mod download

COPY . .
RUN go build -o /pallas-rest rest/rest.go

#
# Prod container
#
FROM gcr.io/distroless/static AS final

LABEL maintainer="Maddosaurus"
USER nonroot:nonroot

COPY --from=builder --chown=nonroot:nonroot /pallas-rest /pallas-rest

EXPOSE 8081
ENTRYPOINT [ "/pallas-rest" ]
