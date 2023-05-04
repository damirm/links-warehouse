FROM golang:1.20-buster AS builder

ADD . /app
RUN cd /app && make compile

FROM alpine
WORKDIR /app

ENV MIGRATIONS_PATH /app/migrations

COPY ./internal/postgres/migrations ${MIGRATIONS_PATH}
COPY --from=builder /app/out/links-warehouse /app/links-warehouse

CMD /app/links-warehouse -config /app/config.yaml
