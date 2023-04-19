FROM golang:1.20-buster AS builder

ADD . /app
RUN cd /app && make compile

FROM alpine
WORKDIR /app

COPY ./internal/postgres/migrations /app/migrations
COPY --from=builder /app/out/links-warehouse /app/links-warehouse

CMD /app/links-warehouse \
    --config /app/config.yaml \
    --db.migrations /app/migrations/ \
    --print-config
