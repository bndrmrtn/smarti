build:
	@go build -o bin/smarti

run: build
	@./bin/smarti run ./language/$(name) --debug

serve: build
	@./bin/smarti serve ./language/$(name)
