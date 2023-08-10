FROM golang:alpine3.18 as builder
WORKDIR /bitmoi-client
COPY . .

RUN go mod tidy
RUN go build

FROM  alpine:latest
USER root
WORKDIR /bitmoi-client
COPY --from=builder /bitmoi-client/bitmoi-client .
COPY ./build ./build

ENTRYPOINT ["./bitmoi-client"]