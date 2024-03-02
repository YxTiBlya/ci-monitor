build:
	@go build -o ./bin/monitor ./cmd/main/main.go

run: build
	@./bin/monitor --cfg config.yaml