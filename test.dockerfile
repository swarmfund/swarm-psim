FROM golang:1.9

WORKDIR /go/src/gitlab.com/swarmfund/psim
COPY . .
ENTRYPOINT ["go", "test"]