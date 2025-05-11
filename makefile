BINARY_NAME=go-ssg
MAIN_PATH=./cmd/main.go

# Default target
all: build

# Build the application
build:
	@echo "Building $(BINARY_NAME)..."
	@go build -o $(BINARY_NAME) $(MAIN_PATH)
	@echo "$(BINARY_NAME) built successfully."

# Run the application
# This will first generate the site and then start the server
run: build
	@echo "Running $(BINARY_NAME)..."
	@./$(BINARY_NAME)

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@go clean
	@rm -f $(BINARY_NAME)
	@rm -rf ./output 
	@echo "Cleaned."

# Create content directory for convenience if it doesn't exist
setup_dirs:
	@mkdir -p content
	@mkdir -p output

.PHONY: all build run clean setup_dirs
