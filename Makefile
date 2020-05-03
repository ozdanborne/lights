build:
	mkdir -p bin/
	go build -o bin/lights .

build-all:
	CGO=enbabled GOOS=darwin go build -o bin/lights-darwin .
	CGO=enbabled GOOS=linux go build -o bin/lights-linux .

run:
	go run .

install:
	ln -s $(PWD)/bin/lights $(HOME)/bin/lights
