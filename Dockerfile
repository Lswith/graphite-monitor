# Docker build
FROM golang
MAINTAINER luke swithenbank swithenbank.luke@gmail.com 

#Installing graphite-monitor
COPY . /go/src/github.com/lswith/graphite-monitor
RUN mkdir /conf
RUN go get -v ./...
RUN go install github.com/lswith/graphite-monitor

#Be sure to put your conf.json into this folder
WORKDIR /conf

ENTRYPOINT ["graphite-monitor"]
VOLUME ["/conf"]

