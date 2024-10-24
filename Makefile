# .PHONY: run


run: build
	@./bin/cryptor

build : template
	@go build -o bin/cryptor

template: css
	@templ generate

css: grammar
	@npm run build

grammar: $(wildcard *.peg)
	@pigeon -o lang/lang.go lang/lang.peg
