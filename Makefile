all: test build

test:
	go test

build:
	go build

clean:
	rm toran
