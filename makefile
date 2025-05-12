BINARY_NAME=go-ssg
MAIN_PATH=./cmd/main.go

# Default target 'all' calls 'build'
all: build

# Build the application: clean, create dirs, tidy modules, then build
build: clean dirs
	@echo "Ensuring dependencies are up to date..."
	@go mod tidy
	@echo "Building $(BINARY_NAME)..."
	@go build -o $(BINARY_NAME) $(MAIN_PATH)
	@echo "$(BINARY_NAME) built successfully."

# Run the application (will show help by default)
run: 
	@echo "Running $(BINARY_NAME)... (use subcommands like generate, serve, edit)"
	@./$(BINARY_NAME)

# Clean build artifacts and output directory
clean:
	@echo "Cleaning..."
	@go clean
	@rm -f $(BINARY_NAME)
	@rm -rf ./output 
	@echo "Cleaned."

# Create content and output directories if they don't exist
dirs:
	@echo "Ensuring content and output directories exist..."
	@mkdir -p content
	@mkdir -p output
	@echo "Directories checked/created."

.PHONY: all build run clean dirs
