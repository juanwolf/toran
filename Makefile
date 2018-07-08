INSTALL_GO_DEP=go get -u github.com/golang/dep/cmd/dep && dep ensure

all: test build

pre:
	$(INSTALL_GO_DEP)

test: pre
	CGO_ENABLED=1 go test -v -coverprofile=coverage.txt -covermode=atomic -race .

build: pre
	go build

clean:
	rm toran
