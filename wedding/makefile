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


.PHONY: migrate migrate-down migrate-new sql-generate
migrate:
	sql-migrate up -config=sql/dbconfig.yml

migrate-local:
	sql-migrate up -config=sql/dbconfig.yml -env="local"

migrate-down:
	sql-migrate down -config=sql/dbconfig.yml

migrate-new:
	sql-migrate new -config=sql/dbconfig.yml $(name)

sql-generate:
	sqlc generate -f sql/sqlc.yaml
