
FROM golang:latest as builder
COPY . $GOPATH/src/TestTask/
WORKDIR $GOPATH/src/TestTask/


RUN go build ./cmd/api/api.go && go build ./cmd/agregator/agregator.go && go build ./cmd/parser/parser.go




FROM golang:latest
RUN apt-get update && apt-get install -y supervisor && mkdir -p $GOPATH/src/TestTask/logs
WORKDIR $GOPATH/src/TestTask/

COPY --from=builder /go/src/TestTask/api /go/src/TestTask/agregator /go/src/TestTask/parser /go/bin/

COPY go_test.sqlite supervisord.conf $GOPATH/src/TestTask/
COPY static $GOPATH/src/TestTask/static/


EXPOSE 8080

CMD supervisord -c $GOPATH/src/TestTask/supervisord.conf
