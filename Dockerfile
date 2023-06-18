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
RUN apk --no-cache add curl
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.16.2/migrate.linux-amd64.tar.gz|tar xvz
RUN mv migrate /usr/bin/migrate
RUN which migrate

COPY ./backend/db/migrate/* ./backend/db/migrate/

COPY ./start_api.sh .
RUN chmod +x start_api.sh

COPY ./wait-for.sh .
RUN chmod +x wait-for.sh