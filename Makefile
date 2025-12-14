lint:
	golangci-lint run

fmt:
	go fmt ./...
	goimports -w .

check: lint fmt

ent-gen:
	ent generate ./internal/ent/schema

.PHONY: lint fmt check ent-gen