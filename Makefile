VERSION=0.2.0-alpha1
BINARY_NAME=zbctl

build:
	@mkdir -p target/bin
	@go build -o target/bin/$(BINARY_NAME) ./cmd/*.go
	@cp cmd/config.toml target/bin/

install:
	@mkdir -p /etc/zeebe/
	@cp target/bin/config.toml /etc/zeebe/
	@cp target/bin/zbctl /usr/local/bin

run:
	@go run cmd/main.go

release:
	@mkdir -p target/bin target/release
	@CGO_ENABLED=0 go build -a -installsuffix cgo -o target/bin/$(BINARY_NAME) cmd/*
	@cp cmd/config.toml target/bin/
	@tar czf $(BINARY_NAME)-$(VERSION).tar.gz target/bin/
	@mv *.tar.gz target/release/
	echo "Success. You can find release at target/release/!"

cov:
	cat coverage_protocol.txt coverage_dumps.txt > coverage.txt
	rm coverage_*.txt

test-client:
	go test zbc/*.go -v

test-protocol:
	go test -race -coverprofile=coverage_protocol.txt -covermode=atomic zbc/zbprotocol/*.go -v

test-hexdump:
	go test -race -coverprofile=coverage_dumps.txt -covermode=atomic tests/test-zbdump/*.go -v

clean:
	@rm -rf ./target *.tar.gz $(BINARY_NAME)
