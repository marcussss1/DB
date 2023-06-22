FROM golang:latest AS build

ADD . /app

WORKDIR /app

RUN CGO_ENABLED=0 go build ./cmd/main/main.go

FROM ubuntu:20.04

COPY . .

ENV TZ=Russia/Moscow
ENV PGVER 12

RUN apt-get -y update && apt-get install -y tzdata
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone
RUN apt-get -y update && apt-get install -y postgresql-$PGVER

USER postgres

RUN /etc/init.d/postgresql start &&\
    psql --command "CREATE USER brabra WITH SUPERUSER PASSWORD 'brabra';" &&\
    createdb -O brabra brabra &&\
    psql -f ./db/db.sql -d brabra &&\
    /etc/init.d/postgresql stop
RUN echo "host all  all    0.0.0.0/0  md5" >> /etc/postgresql/$PGVER/main/pg_hba.conf
RUN echo "listen_addresses='*'" >> /etc/postgresql/$PGVER/main/postgresql.conf

VOLUME  ["/etc/postgresql", "/var/log/postgresql", "/var/lib/postgresql"]

WORKDIR /usr/src/app

COPY . .
COPY --from=build /app/main .

EXPOSE 5000

ENV POSTGRES_USER brabra
ENV POSTGRES_DB brabra
ENV POSTGRES_PASSWORD brabra
ENV POSTGRES_HOST localhost
ENV POSTGRES_PORT 5432
ENV POSTGRES_SSLMODE disable

USER root

CMD service postgresql start && ./main
