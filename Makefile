.PHONY: all run generate clean

all: run

run: generate
	go run ./...

generate:
	go generate ./...

clean:
	go clean
	@echo "Cleaned up build artifacts."
