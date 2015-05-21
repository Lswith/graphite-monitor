FROM golang:1.3
MAINTAINER luke swithenbank swithenbank.luke@gmail.com 

#Installing graphite-monitor
RUN go get github.com/revel/revel
RUN go get github.com/revel/cmd/revel
COPY . /src/github.com/lswith/graphite-monitor
MKDIR /db

WORKDIR /db

ENTRYPOINT ["revel run github.com/lswith/graphite-monitor prod"]
VOLUME ["/db"]