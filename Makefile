build:
	go generate
	go build -o bin/httpthis -v github.com/dainis/httpthis

updatedeps:
	go get -u github.com/kardianos/govendor
	govendor fetch +vendor

initdeps:
	go get -u github.com/kardianos/govendor
	govendor sync

run:
	go generate
	go run *.go

PHONY: build updatedeps initdeps run
