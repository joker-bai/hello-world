FROM golang:1.17.5 AS build-env
ENV GOPROXY https://goproxy.cn
ADD . /go/src/app
WORKDIR /go/src/app
RUN go mod tidy
RUN GOOS=linux GOARCH=386 go build -v -o /go/src/app/app-server

FROM alpine
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories && \
    apk add -U tzdata
RUN ln -sf /usr/share/zoneinfo/Asia/Shanghai  /etc/localtime
COPY --from=build-env /go/src/app/app-server /usr/local/bin/app-server
EXPOSE 8080
CMD [ "app-server" ]