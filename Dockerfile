# 拉取基础镜像
FROM golang

# 设置环境变量 ENV
ENV GOPROXY https://goproxy.cn,direct
ENV GO111MODULE on

# run 执行指令
RUN mkdir /app

# 放置文件，copy 文件会自动解压
ADD . /app/

# cd 指令，当前工作目录
WORKDIR /app
RUN go mod download
RUN go build -o main .
# 规范化 go build
#RUN GOOS=linux GOARCH=amd64 go build ./main.go

# CMD 运行以下指令
CMD ["/app/main"]
