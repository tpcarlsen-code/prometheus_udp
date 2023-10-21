.PHONY: test run build clean

all: test build

pre:
	@mkdir -p build/

build: pre
	GOOS=linux GOARCH=amd64 go build -o build/amd64/prometheus_udp main.go
	GOOS=linux GOARCH=arm64 go build -o build/arm64/prometheus_udp main.go

build_docker:
	docker buildx build --platform linux/amd64 -t promethus_udp:amd64 --load .
	docker buildx build --platform linux/arm64 -t promethus_udp:arm64 --load .

test:
	go test ./...

clean:
	rm -rf build
