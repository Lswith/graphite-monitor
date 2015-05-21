FROM golang:1.4.2
MAINTAINER luke swithenbank swithenbank.luke@gmail.com

#Installing graphite-monitor
RUN go get github.com/revel/revel
RUN go get github.com/revel/cmd/revel
COPY . $GOPATH/src/github.com/lswith/graphite-monitor
RUN mkdir /db

WORKDIR /db

ENTRYPOINT ["revel" ,"run", "github.com/lswith/graphite-monitor","prod"]
CMD ["8090"]
VOLUME ["/db"]
