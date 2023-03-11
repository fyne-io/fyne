# This Makefile is an easy way to run common operations.
# Execute commands this:
# * make test-go
# * make test-integration
#
# Tip: Each command is run on its own line so you can't CD unless you
# connect commands together using operators. See examples:
# A; B    # Run A and then B, regardless of success of A
# A && B  # Run B if and only if A succeeded
# A || B  # Run B if and only if A failed
# A &     # Run A in background.
# Source: https://askubuntu.com/a/539293
#
# Tip: Use $(shell app param) syntax when expanding a shell return value.

.PHONY: test-go
test-go:
	# Run the Go tests.
	go test ./...

.PHONY: test-integration
test-integration:
	# Build the application.
	mkdir -p bin && go build -o bin/goversioninfo cmd/goversioninfo/main.go
	# Test the application.
	PATH="${PATH}:$(shell pwd)/bin" ./testdata/bash/build.sh