all: build
build:
	go build -o ./bin/gim -v
run:  build
	./bin/gim
