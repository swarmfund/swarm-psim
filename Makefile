all: build run

build:
	sh build.sh

run:
	bin/psim run
