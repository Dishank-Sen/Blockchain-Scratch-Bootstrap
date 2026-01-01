build:
	go build -o ./bin/main ./cmd/bootstrap/

run:
	make build
	./bin/main