VERSION=0.1

.PHONY : test deps clean
test: 
	godep go test ./...
deps:
	go get github.com/tools/godep
	godep restore
clean:
	rm clc-*
