.PHONY : test uats
test: 
	go test -v ./clc
uats:
	go test -v ./uats
