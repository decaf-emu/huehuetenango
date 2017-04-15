.PHONY: build run
default: build build_static

build:
	go build -v -o huehuetenango ./cmd/huehuetenango

run: build
	./huehuetenango

build_static:
	cd static && yarn build
