FROM golang:1.13.5-alpine3.10 as build-env
MAINTAINER Xue Bing <xuebing1110@gmail.com>

# repo
RUN cp /etc/apk/repositories /etc/apk/repositories.bak
RUN echo "http://mirrors.aliyun.com/alpine/v3.10/main/" > /etc/apk/repositories
RUN echo "http://mirrors.aliyun.com/alpine/v3.10/community/" >> /etc/apk/repositories

# git
RUN apk update
RUN apk add --no-cache git

# move to GOPATH
RUN mkdir -p /app
WORKDIR /app

# go mod
ENV GOPROXY=https://goproxy.cn,direct
COPY go.mod .
COPY go.sum .
#RUN echo -e "nameserver 10.135.8.110\nnameserver 8.8.8.8" > /etc/resolv.conf & go mod download
RUN go mod download

# build
COPY . .
#COPY etc /app/
RUN go get -u github.com/swaggo/swag/cmd/swag
RUN swag init -g cmd/main.go
RUN go build -o /app/applacation cmd/main.go

## docker image stage
FROM alpine:3.10

# repo
RUN cp /etc/apk/repositories /etc/apk/repositories.bak
RUN echo "http://mirrors.aliyun.com/alpine/v3.10/main/" > /etc/apk/repositories
RUN echo "http://mirrors.aliyun.com/alpine/v3.10/community/" >> /etc/apk/repositories

COPY --from=build-env /app /app


RUN apk update
RUN apk add --upgrade busybox


ENV PORT=8080
EXPOSE 8080
WORKDIR /app
CMD ["/app/applacation"]