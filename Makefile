build-monitor:
	@cd ./monitor && go build -o ../bin/monitor ./cmd/main/main.go

run-monitor: build-monitor
	@./bin/monitor --cfg config.yaml