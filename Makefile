build:
	@go build -o bin/smarti

run: build
	@./bin/smarti run ./language/$(name) --debug
