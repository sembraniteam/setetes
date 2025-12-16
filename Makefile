fmt:
	golangci-lint run --fix

ent-gen:
	ent generate ./internal/ent/schema

schema-apply:
	atlas schema apply --env local

schema-clean:
	atlas schema clean --env local

run:
	go run ./cmd/setetes/main.go start --config config.yml

.PHONY: lint fmt check ent-gen schema-apply schema-clean run