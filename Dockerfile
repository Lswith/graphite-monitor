FROM golang
MAINTAINER luke swithenbank swithenbank.luke@gmail.com 

#Installing graphite-monitor
COPY . /go/src/github.com/lswith/graphite-monitor
RUN mkdir /run
RUN go get -v ./...
RUN go install github.com/lswith/graphite-monitor


WORKDIR /run

ENTRYPOINT ["revel run github.com/lswith/graphite-monitor prod"]
VOLUME ["/run","/conf"]