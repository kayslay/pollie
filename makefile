run: test
	go run cmd/main.go

test: generate
	go test ./...

generate:
	go generate ./...