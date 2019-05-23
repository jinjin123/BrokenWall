FROM golang:latest

WORKDIR /go/src/authcenter

ADD authcenter.go /go/src/authcenter/authcenter.go

RUN go get github.com/gin-gonic/gin

RUN go get gopkg.in/olahol/melody.v1

RUN cd /go/src/authcenter

RUN go build

CMD ["./authcenter"]
