.PHONY: build run clean

build:
	go build -o out/pomodoro ./cmd/pomodoro

run: build
	./out/pomodoro

clean:
	rm -f out/pomodoro

install:
	go install ./cmd/pomodoro

test:
	go test ./...

fmt:
	go fmt ./...

vet:
	go vet ./...
