.PHONY: build run clean

build:
	go build -o out/focus ./cmd/focus

run: build
	./out/focus

clean:
	rm -f out/focus

install:
	go install ./cmd/focus

test:
	go test ./...

fmt:
	go fmt ./...

vet:
	go vet ./...
