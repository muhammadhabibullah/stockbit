
.PHONY: config

config:
	@cp -n config/config.example.json config/config.json

dependency:
	@docker-compose  --profile dependency up -d

tools:
	@docker-compose  --profile tools up -d

build:
	GO111MODULE=on CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o bin/stockbit -mod=vendor -a -installsuffix cgo -ldflags '-w'

swag:
	@echo "> Generate Swagger Docs"
	@if ! command -v swag &> /dev/null; then go install github.com/swaggo/swag/cmd/swag@v1.8.4; fi
	@swag init -o handler/http/docs --ot json
