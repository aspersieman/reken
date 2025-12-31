APP_NAME := reken
BIN_DIR  := bin
PREFIX   := /usr/local
BINDIR   := $(PREFIX)/bin

VERSION  := $(shell git describe --tags --dirty --always 2>/dev/null || echo dev)
LDFLAGS  := -ldflags "-s -w -X main.version=$(VERSION)"

.PHONY: all build install clean run

all: build

build:
	@echo "→ Building $(APP_NAME) ($(VERSION))"
	mkdir -p $(BIN_DIR)
	CGO_ENABLED=1 go build $(LDFLAGS) -o $(BIN_DIR)/$(APP_NAME)

run: build
	./$(BIN_DIR)/$(APP_NAME)

install: build
	@echo "→ Installing to $(BINDIR)"
	install -d $(BINDIR)
	install -m 755 $(BIN_DIR)/$(APP_NAME) $(BINDIR)/$(APP_NAME)

clean:
	rm -rf $(BIN_DIR)

dev:
	go get fyne.io/fyne/v2@latest
	go install fyne.io/tools/cmd/fyne@latest
