# .PHONY: run


run: build
	@./bin/cryptor

build : template
	@go build -o bin/cryptor

template: grammar
	@templ generate

grammar: $(wildcard *.peg)
	@pigeon -o lang/lang.go lang/lang.peg
