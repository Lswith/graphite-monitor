FROM golang:1.3-onbuild
MAINTAINER luke swithenbank swithenbank.luke@gmail.com 

#Installing graphite-monitor
RUN mkdir /db
RUN go get github.com/revel/revel
RUN go get github.com/revel/cmd/revel
RUN go get -v -d -v
RUN go install -v


WORKDIR /db

ENTRYPOINT ["revel run github.com/lswith/graphite-monitor prod"]
VOLUME ["/db","/conf"]