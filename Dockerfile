FROM golang
MAINTAINER luke swithenbank swithenbank.luke@gmail.com 

#Installing graphite-monitor
COPY . /go/src/github.com/lswith/graphite-monitor
RUN mkdir /db
RUN go get github.com/revel/revel
RUN go get github.com/revel/cmd/revel
RUN go get -v ./...


WORKDIR /db

ENTRYPOINT ["revel run github.com/lswith/graphite-monitor prod"]
VOLUME ["/db","/conf"]