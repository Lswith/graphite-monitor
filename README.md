master : [![Build Status](https://travis-ci.org/Lswith/graphite-monitor.svg?branch=master)](https://travis-ci.org/Lswith/graphite-monitor)

develop: [![Build Status](https://travis-ci.org/Lswith/graphite-monitor.svg?branch=develop)](https://travis-ci.org/Lswith/graphite-monitor)
# graphite-monitor
a monitoring daemon for graphite

## Install with Go

If you have go installed on your machine you can run the command:

```sh
go install github.com/lswith/graphite-monitor
```

## Docker Quickstart
```sh
#make the location where the volume will be mounted
mkdir /conf
#edit your configuration
vim /conf/conf.json
#start the container
sudo docker run --name graphite-monitor \
-v /conf:/conf lswith/graphite-monitor
```
