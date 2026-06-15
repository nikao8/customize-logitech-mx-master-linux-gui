BINARY := logid-config-gui
GO      := go
PREFIX  := /usr/local

.PHONY: all build build-release clean install test run

all: build

build: clean
	$(GO) build -tags ci -o $(BINARY)

build-release: clean
	$(GO) build -o $(BINARY)

clean:
	$(RM) $(BINARY)

install: build-release
	install -d $(DESTDIR)$(PREFIX)/bin
	install -m 755 $(BINARY) $(DESTDIR)$(PREFIX)/bin/$(BINARY)

test:
	$(GO) test -tags ci -v ./...

run: build
	./$(BINARY)
