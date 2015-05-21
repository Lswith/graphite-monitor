FROM golang
MAINTAINER luke swithenbank swithenbank.luke@gmail.com 

#Installing graphite-monitor
COPY . /go/src/github.com/lswith/graphite-monitor
RUN mkdir /db
RUN go get -v ./...
RUN go install github.com/lswith/graphite-monitor


WORKDIR /db

ENTRYPOINT ["revel run github.com/lswith/graphite-monitor prod"]
VOLUME ["/db","/conf"]