BIN := ./bin/esa

.PHONY: all
all: clean build

.PHONY: build
build:
	go build -o $(BIN) .

.PHONY: install
install:
	go install .

.PHONY: clean
clean:
	rm -rf $(BIN)
	go clean
