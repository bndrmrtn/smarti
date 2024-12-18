build:
	@go build -o bin/smarti

run: build
	@./bin/smarti

gen: build
	@./bin/rck ./language/$(name)
