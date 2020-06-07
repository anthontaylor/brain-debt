-include secrets.mk

clean:
	go fmt ./...
	go vet ./...

download:
	brew install golang-migrate
	brew install terraform

jaeger:
	docker-compose up -d jaeger

cassandra:
	docker-compose up -d cassandra

build:
	go build -v cmd/main.go cmd/misc.go

run:
	go run cmd/main.go cmd/misc.go -jaeger=localhost:5775 -cassandra=localhost:9042

test:
	go test -v -race ./...

cover:
	go test -race -cover ./... -coverprofile cover.out
	go tool cover -func cover.out
	go tool cover -html=cover.out

docker:
	docker-compose up --build

mock:
	go generate ./...

cqlsh:
	docker exec -it brain-debt_cassandra_1 /opt/dse/bin/cqlsh

cqlsh-dev:
	cqlsh cassandra.us-west-2.amazonaws.com 9142 -u $(cql_dev_user) -p $(cql_dev_password) --ssl

migrate:
	docker exec -i brain-debt_cassandra_1 /opt/dse/bin/cqlsh < ./migrations/keyspace.cql
	migrate -path ./migrations -database cassandra://localhost:9042/brain_debt up
