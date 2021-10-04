FROM golang:1.16-alpine as builder 

COPY go.mod go.sum /go/src/github.com/AjithPanneerselvam/task-etcd/
WORKDIR /go/src/github.com/AjithPanneerselvam/task-etcd 

RUN go mod download
COPY . /go/src/github.com/AjithPanneerselvam/task-etcd/ 

RUN go build -o /task-etcd

EXPOSE 8080

CMD [ "/task-etcd" ]

