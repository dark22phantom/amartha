init:
	@echo "⏳ Initializing project..."
	@echo "📁 Copying example config file..."
ifeq ($(shell test -e config.json && echo "1"),$(shell echo "1"))
	@echo "File config.json already exists"		
else
	@cp example.config.json config.json
endif
	@echo "🏃 Running go mod vendor..."
	@go mod vendor
	@echo "🚀 Ready to run by executing 'make run'"

run:
	@go run cmd/http/main.go

test:
	@echo "🏃 Running tests..."
	@go test -v $$(go list ./... | grep -v /vendor/ | grep -v /cmd/) -cover