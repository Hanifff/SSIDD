FROM golang:latest

COPY . /go/src/app

WORKDIR /go/src/app

RUN mkdir ./fabric

RUN go install

ENTRYPOINT ["/go/bin/ssidd"]
