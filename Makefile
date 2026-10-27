.PHONY: build sqlc

build:
	@go build -o rate-limit-app cmd/api/main.go


sqlc:
	@echo "building sql queries..."
	@sqlc  generate -f internal/clientdb/queries/sqlc.yaml