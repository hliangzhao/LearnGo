FROM golang:latest
WORKDIR /root
ADD ./server.go ./server.go
RUN go build -o server ./server.go
EXPOSE 8080
CMD ["/root/server"]