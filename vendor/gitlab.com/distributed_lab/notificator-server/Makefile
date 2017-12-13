GOPATH=${PWD}/vendor:${PWD}
LINTER_LOG=${PWD}/linter.log

lint:
	gometalinter ./src/... | grep -v "should have comment" > ${LINTER_LOG} || true
	[ ! -s ${LINTER_LOG} ] || cat ${LINTER_LOG}

build:
	./build