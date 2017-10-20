VERSION=0.2.0
BINARY_NAME=zbctl
ZBC_PATH=$(GOPATH)/src/github.com/zeebe-io/zbc-go
PREFIX=/usr/local

build:
	@mkdir -p target/bin
	@go build -o target/bin/$(BINARY_NAME) ./cmd/*.go
	@cp cmd/config.toml target/bin/

.PHONY: install
install:
	mkdir -p $(PREFIX)/zeebe
	install -m 644 target/bin/config.toml $(PREFIX)/zeebe/
	install -m 755 target/bin/zbctl $(PREFIX)/zeebe/
	ln -sf $(PREFIX)/zeebe/zbctl $(PREFIX)/bin/zbctl

.PHONY: uninstall
uninstall:
	rm -rf $(PREFIX)/bin/zbctl
	rm -rf $(PREFIX)/zeebe

run:
	@go run cmd/main.go

cov:
	cat coverage_protocol.txt coverage_dumps.txt > coverage.txt
	rm coverage_*.txt

test-client:
	go test zbc/*.go -v

test-protocol:
	go test -race -coverprofile=coverage_protocol.txt -covermode=atomic zbc/zbprotocol/*.go -v

test-hexdump:
	go test -race -coverprofile=coverage_dumps.txt -covermode=atomic tests/test-zbdump/*.go -v

test-integration:
	go test -race -coverprofile=coverage_integration.txt -covermode=atomic tests/test-broker/*.go -v

test-all: test-client test-protocol test-hexdump test-integration

clean:
	@rm -rf ./target *.tar.gz $(BINARY_NAME)
