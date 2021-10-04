FROM golang:1.16-alpine as builder 

COPY go.mod go.sum /go/src/github.com/AjithPanneerselvam/todo/
WORKDIR /go/src/github.com/AjithPanneerselvam/todo 

RUN go mod download
COPY . /go/src/github.com/AjithPanneerselvam/todo/ 

RUN go build -o /todo

EXPOSE 8080

CMD [ "/todo" ]

