FROM ubuntu:14.04

ENV GOPATH /go

RUN apt-get update && apt-get install -y stockfish && rm -rf /var/lib/apt/lists/*

ADD map/map /bin/map

EXPOSE 80
