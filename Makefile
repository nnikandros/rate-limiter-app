.PHONY: build sqlc

build:
	@go build -o rate-limit-app cmd/api/main.go

build-static:
	@CGO_ENABLED=0 GOOS=linux go build -o rate-limit-app cmd/api/main.go

build-arm:
	@echo "Cross compiling for arm64..."
	@CGO_ENABLED=0 GOOS=linux GOARCH=arm64 CC=aarch64-linux-gnu-gcc go build -o rate-limit-app cmd/api/main.go


sqlc:
	@echo "building sql queries..."
	@sqlc  generate -f internal/clientdb/queries/sqlc.yaml