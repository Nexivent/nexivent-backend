.PHONY: all run build clean

all: build

run: build
	@./bin/api

build:
	@go build -o bin/api ./cmd/api

clean:
	@rm -rf bin/api