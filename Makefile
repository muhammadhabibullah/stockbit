
.PHONY: config

config:
	@cp -n config/config.example.json config/config.json

dependency:
	@docker-compose  --profile dependency up -d

build:
	GO111MODULE=on CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o bin/stockbit -mod=vendor -a -installsuffix cgo -ldflags '-w'
