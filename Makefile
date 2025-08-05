build:
	@go build -o bin/simple-mock-api.exe main.go

test:
	@go test ./...

testfast:
	@go test ./... -failfast

testv:
	@go test -v ./...

run: build
	@./bin/simple-mock-api.exe