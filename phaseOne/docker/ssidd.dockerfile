FROM golang

COPY . /go/src/app

WORKDIR /go/src/app

RUN ls .

RUN go install -v ./app/ssidd

ENTRYPOINT ["/go/bin/ssidd"]
