build:
	go build -o bin/lib

run: build
	./bin/lib
