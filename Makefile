.PHONY: test
test:
	@go test ./...

.PHONY: tparse
tparse:
	@go test -race -count=1 ./... -json -cover -coverpkg=./... | tparse

.PHONY: coverage
coverage:
	@go test -cover -coverprofile=coverage.out ./...
	@go tool cover -func=coverage.out

.PHONY: lint
lint:
	@golangci-lint run ./...
	@govulncheck ./...

release:
	@goreleaser --rm-dist

tidy:
	@go mod tidy

deps:
	@go install golang.org/x/vuln/cmd/govulncheck@latest
	@go install github.com/mfridman/tparse@latest
