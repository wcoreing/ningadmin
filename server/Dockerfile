FROM golang:1.20.13-alpine3.18
RUN go env -w GO111MODULE=on && \
go env -w GOPROXY=https://goproxy.cn,direct;
WORKDIR /www
COPY ./ ./
CMD ["go" ,"build", "main.go"]
