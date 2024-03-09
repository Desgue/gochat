run: build
@./bin/gochat

build: 
@go build -o bin/gochat

test:
@go test -v ./...