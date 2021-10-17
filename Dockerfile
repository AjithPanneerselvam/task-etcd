FROM golang:1.16-alpine 

WORKDIR /go/src/github.com/AjithPanneerselvam/task-etcd 

COPY go.mod go.sum /go/src/github.com/AjithPanneerselvam/task-etcd/
RUN go mod download

COPY . /go/src/github.com/AjithPanneerselvam/task-etcd/ 
RUN go build -o /task-etcd

COPY /static/index.html /static/

EXPOSE 8080

# Command to run the executable
ENTRYPOINT ["/task-etcd"]
