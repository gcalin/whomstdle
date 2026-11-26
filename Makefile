.PHONY: check fmt cover

# Run every check. GitHub runs this same command on each push.
# Prints only short step names, plus details when something fails.
check:
	@echo "==> format"
	@out=$$(gofmt -l $$(go list -f '{{.Dir}}' ./...)); \
		if [ -n "$$out" ]; then echo "Not formatted (fix with: make fmt):"; echo "$$out"; exit 1; fi
	@echo "==> go.mod tidy"
	@go mod tidy -diff
	@echo "==> lint (includes go vet)"
	@golangci-lint run ./...
	@echo "==> test"
	@go test -race -cover ./...
	@echo "All checks passed."

# Format all Go files in place.
fmt:
	go fmt ./...

# Coverage report: prints the total, and writes coverage.html (open it to see which lines tests run).
cover:
	@go test -coverprofile=coverage.out ./... > /dev/null
	@go tool cover -func=coverage.out | tail -1
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Wrote coverage.html"
