build:
	mkdir -p bin/
	go build -o bin/lights .

run:
	go run .

install:
	ln -s $(PWD)/bin/lights $(HOME)/bin/lights
