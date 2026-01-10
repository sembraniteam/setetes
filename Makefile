VERSION := $(shell cat ./VERSION)

fmt:
	golangci-lint run --fix

ent-gen:
	ent generate ./internal/ent/schema

schema-apply:
	atlas schema apply --env local --auto-approve

schema-clean:
	atlas schema clean --env local --auto-approve

run:
	go run ./cmd/setetes/main.go start --config config.yml

seed:
	go run ./cmd/setetes/main.go seed --config config.yml

build:
	GOOS=linux CGO_ENABLED=0 go build -v -ldflags='-s -w -X main.version=$(VERSION)' -o /bin/app cmd/setetes/main.go
	chmod +x /bin/app

.PHONY: fmt ent-gen schema-apply schema-clean run seed build