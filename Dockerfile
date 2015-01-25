FROM golang:1.4
MAINTAINER Andreas Koch <andy@ak7.io>

# Build
ADD . /go
RUN go run make.go -install

EXPOSE 7070

CMD ["/go/bin/postdeploy", ":7070", "-config", "postdeploy.conf.js"]