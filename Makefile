# =========================
# Makefile — bootstrap
# =========================

APP_NAME=bootstrap
CMD_PATH=./cmd/bootstrap
BIN_DIR=bin

GOOS ?= $(shell go env GOOS)

ifeq ($(GOOS),windows)
	EXT=.exe
else
	EXT=
endif

BIN=$(BIN_DIR)/$(APP_NAME)$(EXT)

.PHONY: build run clean

# -------------------------
# Build
# -------------------------
build:
	go build -o "$(BIN)" "$(CMD_PATH)"

# -------------------------
# Run (foreground)
# -------------------------
run: build
	$(BIN)

# -------------------------
# Clean
# -------------------------
clean:
	rm -rf "$(BIN_DIR)"
