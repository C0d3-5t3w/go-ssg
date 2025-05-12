BINARY_NAME=go-ssg
MAIN_PATH=./cmd/main.go

all: build

build: clean dirs
	@echo "Ensuring dependencies are up to date..."
	@go mod tidy
	@echo "Building $(BINARY_NAME)..."
	@go build -o $(BINARY_NAME) $(MAIN_PATH)
	@echo "$(BINARY_NAME) built successfully."

run: 
	@echo "Running $(BINARY_NAME)... (use subcommands like generate, serve, edit)"
	@./$(BINARY_NAME)

clean:
	@echo "Cleaning..."
	@go clean
	@rm -f $(BINARY_NAME)
	@rm -rf ./output 
	@echo "Cleaned."

dirs:
	@echo "Ensuring content and output directories exist..."
	@mkdir -p content
	@mkdir -p output
	@echo "Directories checked/created."

.PHONY: all build run clean dirs
