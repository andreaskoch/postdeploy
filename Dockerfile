FROM golang:alpine
MAINTAINER Andreas Koch <andy@ak7.io>

# Add code
ADD . /go/src/github.com/andreaskoch/postdeploy

# Build
RUN cd /go/src/github.com/andreaskoch/postdeploy && \
    go build -o /bin/postdeploy && \
    rm -rf /go/pkg

# Config
RUN mkdir -p /etc/postdeploy/conf
ADD conf/ping-sample.json /etc/postdeploy/conf/postdeploy.json

EXPOSE 7070

CMD ["postdeploy", "-binding=:7070", "-config=/etc/postdeploy/conf/postdeploy.json"]
