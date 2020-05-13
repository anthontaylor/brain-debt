clean:
	go fmt ./...
	go vet ./...

run:
	docker-compose up -d jaeger
	go run cmd/main.go cmd/misc.go -jaeger=localhost:5775

test:
	go test -race ./...

cover:
	go test -race -cover ./... -coverprofile cover.out
	go tool cover -func cover.out
	go tool cover -html=cover.out

docker:
	docker-compose up --build
