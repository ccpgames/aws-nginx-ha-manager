all: provision test build

provision:
	go get -t ./...

test:
	go test -v ./...

build:
	mkdir -p build
	go build -o build/aws-nginx-ha-manager
