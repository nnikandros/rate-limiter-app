.PHONY: build sqlc

build:
	@go build -o rate-limit-app cmd/api/main.go

build-static:
	@CGO_ENABLED=0 GOOS=linux go build -o rate-limiter cmd/api/main.go

build-arm:
	@echo "Cross compiling for arm64..."
	@CGO_ENABLED=0 GOOS=linux GOARCH=arm64 CC=aarch64-linux-gnu-gcc go build -o rate-limiter cmd/api/main.go

sync: build-arm
	@echo "sync"
	@rsync -av rate-limiter mock_db.json  hetzner-app-runner:/home/app-runner/applications/rate-limiter

reload:
	@echo "reloading service"
	@echo ${APP_RUNNER_PASSWORD} | ssh hetzner-app-runner -t sudo -S systemctl reload rate-limiter.service 2> /dev/null


status:
	@ssh hetzner-app-runner -t systemctl status rate-limiter.service 2> /dev/null

journal:
	@ssh hetzner-app-runner -t journalctl -fu rate-limiter.service