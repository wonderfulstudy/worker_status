FROM golang:1.23.4 As builder
WORKDIR /go/src
ADD . /go/src
RUN cd /go/src && \
    go version && \
    go mod download && \
    go mod tidy && \
    go build -o nyx -v .
 
# FROM alpine:3.21.2
# WORKDIR /app
# COPY --from=builder /go/src/nyx /app/nyx
# COPY ./config.yaml /app/config.yaml
# ENTRYPOINT ./nyx