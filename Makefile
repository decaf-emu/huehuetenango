.PHONY: deps generate build run
default: deps generate build

deps:
	go get -u github.com/mailru/easyjson/...
	$(MAKE) -C static deps

generate:
	easyjson -all -pkg pkg/models
	easyjson -all -pkg pkg/importer/schema

build:
	go build -v -o huehuetenango ./cmd/huehuetenango
	$(MAKE) -C static build

run: generate build
	./huehuetenango
