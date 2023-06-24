FROM golang:alpine3.18 as builder
WORKDIR /bitmoi
COPY . .

RUN go mod tidy
RUN go build

FROM  alpine:latest
USER root
WORKDIR /bitmoi
COPY --from=builder /bitmoi/bitmoi .
COPY ./a.env .
COPY ./Makefile .

RUN apk update
RUN apk add make