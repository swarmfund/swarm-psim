FROM golang:1.9-alpine

WORKDIR /go/src/gitlab.com/swarmfund/psim
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /binary -v gitlab.com/swarmfund/psim/cmd/psim

FROM ubuntu:latest
COPY --from=0 /binary .
RUN apk update \
 && apk add ca-certificates \
 && rm -rf /var/cache/apk/*
ENTRYPOINT ["./binary", "--config", "/config.yaml"]
