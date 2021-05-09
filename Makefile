BIN := esa

.PHONY: all
all: clean build

.PHONY: build
build:
	go build -o $(BIN) .

.PHONY: clean
clean:
	rm -rf $(BIN)
	go clean
