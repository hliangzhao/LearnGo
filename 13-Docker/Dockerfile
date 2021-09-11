#这种写法非常不推荐，因为每一行都要构建新的一层！最后生成的image庞大无比（约1G）
FROM golang:latest

RUN mkdir /app

WORKDIR /app

ADD ./server.go /app

RUN go build -o main ./server.go

EXPOSE 8080

CMD /app/main