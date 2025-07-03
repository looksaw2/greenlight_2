buildHTTP:
	@echo "start to build the HTTP Server"
	@go build -o ./bin ./cmd/main/main.go
runHTTP: buildHTTP
	@echo "start to run this HTTP Server"
	@./bin/main