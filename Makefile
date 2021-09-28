default: build

build:
	go build

install:
	go install

test:
	go test --race --coverprofile=coverage.txt ./...

coverage:
	go tool cover --html=coverage.txt