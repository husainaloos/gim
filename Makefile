all: build
build:
	go build -o ./bin/gim -v
run: build
	./bin/gim
try: build
	./bin/gim ./tmp/sample.txt
