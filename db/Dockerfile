FROM postgres:14
ENV POSTGRES_PASSWORD postgres
ENV POSTGRES_DB testdb
COPY init.sql /docker-entrypoint-initdb.d/
