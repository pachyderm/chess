FROM ubuntu:14.04

ENV GOPATH /go

RUN apt-get update && apt-get install -y software-properties-common python python-pip && rm -rf /var/lib/apt/lists/*
RUN add-apt-repository "deb http://us.archive.ubuntu.com/ubuntu/ trusty universe multiverse"
RUN apt-get update && apt-get install -y crafty crafty-books-medium && rm -rf /var/lib/apt/lists/*
RUN cp /usr/share/doc/crafty/setup_crafty.sh . && chmod +x setup_crafty.sh && ./setup_crafty.sh && rm setup_crafty.sh
RUN pip install pgnparser

ADD map/map /bin/map
ADD map/crafty-to-json /bin/crafty-to-json
ADD reduce/reduce /bin/reduce

EXPOSE 80
CMD /bin/crafty-server
