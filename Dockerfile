FROM golang

ARG GIN_MODE_C=debug
ENV GIN_MODE=$GIN_MODE_C

WORKDIR /go/src/github.com/WodBoard/wod-api
ADD . /go/src/github.com/WodBoard/wod-api

RUN go get -u github.com/cespare/reflex
RUN go get -u github.com/golang/dep/cmd/dep && \
    dep ensure
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./wod-api

CMD reflex -r '\.go$$' -s -- sh -c 'CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./wod-api && ./wod-api'