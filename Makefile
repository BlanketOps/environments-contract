BUF ?= buf

.PHONY: all
all: lint generate

.PHONY: generate
generate:
	$(BUF) generate

.PHONY: lint
lint:
	$(BUF) lint

.PHONY: format
format:
	$(BUF) format -w

.PHONY: clean
clean:
	find blanketops -name "*.pb.go" -delete

.PHONY: regen
regen: clean generate
