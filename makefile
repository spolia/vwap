test-unit:
	@go test -race ./internal/vwapcalculator
	@go test -race ./internal/websocket/coinbase

test-integration:
	@go test -race ./internal

upgrade:
	@echo "Upgrading dependencies..."
	@go get -u
	@go mod tidy
	
run:
	@go run main.go

build:
	@go build -o vwap main.go

clean:
	@rm -rf wvap