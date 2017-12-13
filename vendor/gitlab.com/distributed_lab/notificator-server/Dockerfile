FROM golang:1.9

RUN mkdir -p $GOPATH/src/gitlab.com/distributed_lab/notificator-server
COPY . $GOPATH/src/gitlab.com/distributed_lab/notificator-server
WORKDIR $GOPATH/src/gitlab.com/distributed_lab/notificator-server

RUN ./docker-setup.sh

WORKDIR /app