FROM golang:1.4
MAINTAINER Andreas Koch <andy@ak7.io>

# Build
ADD . /go
RUN go run make.go -install

# Config
RUN mkdir -p /etc/postdeploy/conf
ADD conf/ping-sample.json /etc/postdeploy/conf/postdeploy.json

EXPOSE 7070

CMD ["/go/bin/postdeploy", "-binding=:7070", "-config=/etc/postdeploy/conf/postdeploy.json"]