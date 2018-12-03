GO_FILES = main.go $(wildcard */*.go) go.mod go.sum

.PHONY: all
all: dbp

dbp: $(GO_FILES)
	go build -o $@

.PHONY: run
run: dbp
	./$<

.PHONY: clean
clean:
	rm -rf dbp
