FROM golang:1.9.2-alpine3.6
MAINTAINER Xue Bing <xuebing1110@gmail.com>

# repo
RUN cp /etc/apk/repositories /etc/apk/repositories.bak
RUN echo "http://mirrors.aliyun.com/alpine/v3.6/main/" > /etc/apk/repositories
RUN echo "http://mirrors.aliyun.com/alpine/v3.6/community/" >> /etc/apk/repositories

# timezone
RUN apk update
RUN apk add --no-cache tzdata \
    && echo "Asia/Shanghai" > /etc/timezone \
    && ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime

# move to GOPATH
RUN mkdir -p /go/src/github.com/pingfen/noticeplat-server
COPY . $GOPATH/src/github.com/pingfen/noticeplat-server
WORKDIR $GOPATH/src/github.com/pingfen/noticeplat-server/

# build
RUN mkdir -p /app
RUN go build -o /app/msgpack cmd/main.go

WORKDIR /app
EXPOSE 8080
ENV GIN_MODE=release
CMD ["/app/msgpack"]