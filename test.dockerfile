FROM golang:1.9-alpine

WORKDIR /go/src/gitlab.com/swarmfund/psim
COPY . .
ENTRYPOINT ["go", "test"]