FROM ubuntu:14.04

ENV GOPATH /go

RUN apt-get update && apt-get install -y stockfish curl && rm -rf /var/lib/apt/lists/*
RUN \
  curl -sSL https://storage.googleapis.com/golang/go1.6.linux-amd64.tar.gz | tar -C /usr/local -xz && \
  mkdir -p /go/bin
ENV PATH /go/bin:/usr/local/go/bin:$PATH
ENV GOPATH /go

ADD . /go/src/annotate
RUN go install annotate
