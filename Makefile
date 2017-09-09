.PHONY: deps build run
default: deps build

deps:
	$(MAKE) -C static deps

build:
	go build -v -o huehuetenango ./cmd/huehuetenango
	$(MAKE) -C static build

run: generate build
	./huehuetenango
