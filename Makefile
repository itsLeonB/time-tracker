.PHONY: api-hotreload lint

# Run the backend API with hot-reload, make sure you installed `air`
api-hotreload:
	air --build.cmd "go build -o bin/app cmd/app/main.go" --build.bin "./bin/app"

lint:
	golangci-lint run ./...
