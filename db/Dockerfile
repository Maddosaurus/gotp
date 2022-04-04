# syntax=docker/dockerfile:1

# This Dockerfile is optional, you don't need this if you use the included docker-compose.yml
FROM postgres:14
ENV POSTGRES_PASSWORD passpass
ENV POSTGRES_DB gotp
COPY init.sql /docker-entrypoint-initdb.d/
