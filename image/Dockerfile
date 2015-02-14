FROM ubuntu:14.04

ENV GOPATH /go

RUN apt-get update && apt-get install -y software-properties-common stockfish && rm -rf /var/lib/apt/lists/*
RUN add-apt-repository "deb http://us.archive.ubuntu.com/ubuntu/ trusty universe multiverse"

ADD map/map /bin/map

EXPOSE 80
