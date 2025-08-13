templ:
	templ generate --watch
	@echo "Template generation started. Watching for changes..."
	@echo "To stop the process, use Ctrl+C."

.PHONY: templ
templ: templ

build:
	go build -o myhometask ./cmd/server/main/main.go


.PHONY: build
build: build

.PHONY: tools
tools:
	@go build -o ./tools ./cmd/tools/

.PHONY: sqlc
sqlc:
	sqlc -f store/sqlc.yaml generate

