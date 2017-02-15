RDIR = run

build:
	go install -v ./...

# Run tests
test:
	go test ./...

# Run locally for testing
run: build
	rm -rf $(RDIR)
	mkdir -p $(RDIR)/bak
	ssbd
