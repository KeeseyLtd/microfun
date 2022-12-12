.PHONY:  dependencies test test test build run local local-debug debug swagger

dependencies:
	go install github.com/githubnemo/CompileDaemon@latest

test:
	go test ./...

build:
	go build ./...

generate:
	bash ./scripts/mocks.sh

run:
	CompileDaemon -polling -polling-interval 1000 -build="go build main.go" -command=./main

debug:
	CompileDaemon -polling -build="go build main.go" -command="dlv --listen=:2345 --headless=true --log=true --accept-multiclient --api-version=2 exec ./main"

local:
	APP_CMD=run docker compose -f build/local/docker-compose.yml up

local-debug:
	APP_CMD=debug docker compose -f build/local/docker-compose.yml up

swagger:
	swagger generate spec -o api/swagger.yaml && sed -i "s/default/x-example/g" api/swagger.yaml
