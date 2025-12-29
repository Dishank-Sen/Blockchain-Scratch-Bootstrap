build:
	go build -o ./bin/main ./cmd/bootstrap/main.go

run:
	make build
	./bin/main