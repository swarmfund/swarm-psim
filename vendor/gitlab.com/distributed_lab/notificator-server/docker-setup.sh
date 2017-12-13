#!/usr/bin/env sh

prepare() {
    eval $(ssh-agent -s) && \
    chmod 400 .id_rsa && \
    ssh-add .id_rsa && \
    mkdir -p ~/.ssh && \
    git config --global url.ssh://git@gitlab.com/.insteadOf https://gitlab.com/ && \
    echo "Host *\n\tStrictHostKeyChecking no\n\n" > ~/.ssh/config && \
    go get -u github.com/golang/dep/cmd/dep
}

copy_files() {
    mkdir /app && \
    ls -la && \
    cp -rf ./run /app/ && \
    cp -rf ./bin /app/ && \
    cp -rf ./conf.yaml /app/ && \
    cp -rf ./migrations /app/
}

prepare && \
echo "dep ensure ..." && dep ensure && \
./run build && \
copy_files
