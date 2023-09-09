FROM golang
ENV TZ=Asia/Shanghai
WORKDIR /app
ADD . /app
ENV GOPROXY=https://mirrors.aliyun.com/goproxy/
RUN go mod download
RUN go build main.go
ENTRYPOINT ["./main"]