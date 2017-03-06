FROM golang:1.7-alpine

RUN mkdir -p /go/src/github.com/dollarshaveclub/polymerase

ADD . /go/src/github.com/dollarshaveclub/polymerase

WORKDIR /go/src/github.com/dollarshaveclub/polymerase

RUN go install

CMD ["/go/bin/polymerase"]