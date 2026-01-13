FROM golang

ADD kongExample /go/src/github.com/username/kongExample

RUN go get github.com/gorilla/mux

RUN go install github.com/username/kongExample

ENTRYPOINT /go/bin/kongExample