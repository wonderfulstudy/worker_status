FROM golang:1.23.4-alpine AS builder

WORKDIR /app
COPY ./ /app

RUN go build -o /app/blackbox_exporter -v .

FROM alpine:latest
LABEL maintainer="The Prometheus Authors <prometheus-developers@googlegroups.com>"

COPY --from=builder /app/blackbox_exporter /bin/blackbox_exporter
COPY blackbox.yml       /etc/blackbox_exporter/config.yml

EXPOSE      9115
ENTRYPOINT  [ "/bin/blackbox_exporter" ]
CMD         [ "--config.file=/etc/blackbox_exporter/config.yml" ]