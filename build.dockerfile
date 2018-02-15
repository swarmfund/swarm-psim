FROM golang:1.9

WORKDIR /go/src/gitlab.com/swarmfund/psim
COPY . .
RUN GOOS=linux go build -ldflags "-s" -o /binary gitlab.com/swarmfund/psim/psim/cmd/psim

# can't use alpine because of ethereum cgo extensions
FROM ubuntu:latest
COPY --from=0 /binary .
RUN apt-get update \
 && apt-get install -y ca-certificates
ENTRYPOINT ["./binary", "--config", "/config.yaml"]
