FROM golang:1.10.1

ENV GOPATH /go

ENV PATH /go/bin:$PATH
RUN go get github.com/tools/godep

RUN go get github.com/jteeuwen/go-bindata/...

RUN mkdir -p /go/src/github.com/votinginfoproject/sms-worker
WORKDIR /go/src/github.com/votinginfoproject/sms-worker

RUN touch .env

COPY . /go/src/github.com/votinginfoproject/sms-worker

RUN go-bindata -prefix "data" -pkg "data" -o data/data.go data/raw

RUN godep go test ./...

RUN godep go build -o sms-worker sms-worker.go

CMD ./sms-worker
