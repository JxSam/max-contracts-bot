.DEFAULT_GOAL = local

local:
	CONFIG_PATH=config/local.yaml go run cmd/main.go

prod:
	CONFIG_PATH=config/prod.yaml go run cmd/main.go

.PHONY: deploy
deploy:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -ldflags '-extldflags "-static"' -o deploy/max-contracts-bot cmd/main.go