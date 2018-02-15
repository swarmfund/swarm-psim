export GOCONVEY_REPORTER=story
go test $(go list ./... | grep -v /vendor/) -parallel 1 | egrep -v '\?.+ [no test files]'
