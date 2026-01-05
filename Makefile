APP_NAME := calc
APP_TITLE := Reken

BIN_DIR := bin
PREFIX := $(HOME)/.local

BINDIR  := $(PREFIX)/bin
APPDIR  := $(PREFIX)/share/applications
ICONDIR := $(PREFIX)/share/icons/hicolor/256x256/apps

DESKTOP_FILE := $(APP_NAME).desktop
ICON_NAME := $(APP_NAME)

VERSION  := $(shell git describe --tags --dirty --always 2>/dev/null || echo dev)
LDFLAGS  := -ldflags "-s -w -X main.version=$(VERSION)"

.PHONY: all build install uninstall clean run

all: build

build:
	@echo "→ Building $(APP_TITLE) ($(VERSION)): $(BIN_DIR)/$(APP_NAME)"
	mkdir -p $(BIN_DIR)
	CGO_ENABLED=1 go build $(LDFLAGS) -o $(BIN_DIR)/$(APP_NAME)

run: build
	./$(BIN_DIR)/$(APP_NAME)

install: build
	@echo "→ Installing binary"
	install -Dm755 $(BIN_DIR)/$(APP_NAME) $(BINDIR)/$(APP_NAME)

	@echo "→ Installing icon"
	install -Dm644 assets/icon.png $(ICONDIR)/$(ICON_NAME).png

	@echo "→ Installing desktop entry"
	@install -d $(APPDIR)
	@printf '%s\n' \
		'[Desktop Entry]' \
		'Type=Application' \
		'Name=$(APP_TITLE)' \
		'Exec=$(APP_NAME)' \
		'Icon=$(ICON_NAME)' \
		'Categories=Utility;Calculator;' \
		'StartupNotify=true' \
		> $(APPDIR)/$(DESKTOP_FILE)

	@echo "→ Updating desktop database (if available)"
	@-update-desktop-database $(APPDIR) >/dev/null 2>&1 || true

	@echo "✓ $(APP_TITLE) installed"

uninstall:
	@echo "→ Removing binary"
	rm -f $(BINDIR)/$(APP_NAME)

	@echo "→ Removing icon"
	rm -f $(ICONDIR)/$(ICON_NAME).png

	@echo "→ Removing desktop entry"
	rm -f $(APPDIR)/$(DESKTOP_FILE)

	@echo "→ Updating desktop database (if available)"
	-update-desktop-database $(APPDIR) >/dev/null 2>&1 || true

	@echo "✓ $(APP_TITLE) uninstalled"

clean:
	rm -rf $(BIN_DIR)
