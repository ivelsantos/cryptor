# .PHONY: run


run: build
	@./bin/cryptor

build : grammar
	@go build -o bin/cryptor

grammar: $(wildcard *.peg)
	@pigeon -o lang/lang.go lang/lang.peg
