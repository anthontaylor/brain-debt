clean:
	go fmt ./...
	go vet ./...

run:
	go run cmd/main.go

test:
	go test -race ./...

cover:
	go test -race -cover ./... -coverprofile cover.out
	go tool cover -func cover.out
	go tool cover -html=cover.out
