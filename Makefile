BINARY   := logid-config-gui
GO       := go
PREFIX   := /usr/local
VERSION  := $(shell git describe --tags --always 2>/dev/null || echo "dev")
RELEASE  := release

.PHONY: all build build clean install test run
.PHONY: build-debian build-arch build-fedora build-all

all: release run


build: clean
	$(GO) build -o $(BINARY) ./src/cmd/main.go

build-debian: clean
	$(GO) build -o $(BINARY) ./src/cmd/main.go
	mkdir -p $(RELEASE)
	TMP=$$(mktemp -d) && \
	cp $(BINARY) $$TMP/ && \
	cp install/install_debian.sh $$TMP/install.sh && \
	chmod +x $$TMP/install.sh && \
	tar -czf $(RELEASE)/$(BINARY)-debian-$(VERSION).tar.gz -C $$TMP . && \
	$(RM) -r $$TMP $(BINARY) && \
	echo "==> Release criado: $(RELEASE)/$(BINARY)-debian-$(VERSION).tar.gz"

build-arch: clean
	$(GO) build -o $(BINARY) ./src/cmd/main.go
	mkdir -p $(RELEASE)
	TMP=$$(mktemp -d) && \
	cp $(BINARY) $$TMP/ && \
	cp install/install_arch.sh $$TMP/install.sh && \
	chmod +x $$TMP/install.sh && \
	tar -czf $(RELEASE)/$(BINARY)-arch-$(VERSION).tar.gz -C $$TMP . && \
	$(RM) -r $$TMP $(BINARY) && \
	echo "==> Release criado: $(RELEASE)/$(BINARY)-arch-$(VERSION).tar.gz"

build-fedora: clean
	$(GO) build -o $(BINARY) ./src/cmd/main.go
	mkdir -p $(RELEASE)
	TMP=$$(mktemp -d) && \
	cp $(BINARY) $$TMP/ && \
	cp install/install_fedora.sh $$TMP/install.sh && \
	chmod +x $$TMP/install.sh && \
	tar -czf $(RELEASE)/$(BINARY)-fedora-$(VERSION).tar.gz -C $$TMP . && \
	$(RM) -r $$TMP $(BINARY) && \
	echo "==> Release criado: $(RELEASE)/$(BINARY)-fedora-$(VERSION).tar.gz"

release: build-debian build-arch build-fedora

clean:
	$(RM) $(BINARY)
	$(RM) -r $(RELEASE)

install: build
	install -d $(DESTDIR)$(PREFIX)/bin
	install -m 755 $(BINARY) $(DESTDIR)$(PREFIX)/bin/$(BINARY)

run: build
	sudo ./$(BINARY)

