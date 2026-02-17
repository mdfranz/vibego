# Makefile for vibego

# Variables
GO := go
BIN_DIR := bin
CMD_DIR := cmd
# Get all subdirectories in cmd/ as the list of binaries to build
COMMANDS := $(shell [ -d $(CMD_DIR) ] && ls $(CMD_DIR))
BINARIES := $(addprefix $(BIN_DIR)/,$(COMMANDS))

# Phony targets
.PHONY: all build test clean fmt vet lint tidy help

# Default target
all: clean fmt vet test build

# Build all binaries dynamically
build: $(BINARIES)

$(BIN_DIR)/%: $(CMD_DIR)/%/main.go
	@echo "Building $@..."
	@mkdir -p $(BIN_DIR)
	$(GO) build -o $@ ./$(CMD_DIR)/$*

# Run tests
test:
	@echo "Running tests..."
	$(GO) test ./... -v

# Clean build artifacts
clean:
	@echo "Cleaning up..."
	@rm -rf $(BIN_DIR)

# Format code
fmt:
	@echo "Formatting code..."
	$(GO) fmt ./...

# Vet code
vet:
	@echo "Vetting code..."
	$(GO) vet ./...

# Lint code (checks if golangci-lint is installed)
lint:
	@echo "Linting code..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not installed. Skipping. (See TL-1 in GEMINI.md)"; \
	fi

# Tidy modules
tidy:
	@echo "Tidying modules..."
	$(GO) mod tidy

# Help target
help:
	@echo "Available targets:"
	@echo "  all    - Clean, fmt, vet, test, and build (default)"
	@echo "  build  - Build all binaries in $(CMD_DIR)/"
	@echo "  test   - Run all tests"
	@echo "  clean  - Remove $(BIN_DIR)/"
	@echo "  fmt    - Format code using 'go fmt'"
	@echo "  vet    - Vet code using 'go vet'"
	@echo "  lint   - Run golangci-lint (if available)"
	@echo "  tidy   - Run 'go mod tidy'"
