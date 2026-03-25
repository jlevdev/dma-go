goose: ## make goose command="up-to version-id"
	goose -dir ./migrations postgres "host=localhost user=dma password=dma dbname=dma sslmode=disable" $(command)

migrate-up: ## Migrate up
	goose -dir ./migrations postgres "host=localhost user=dma password=dma dbname=dma sslmode=disable" up

migrate-down: ## Migrate down
	goose -dir ./migrations postgres "host=localhost user=dma password=dma dbname=dma sslmode=disable" down

help: ## Display this help message
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n\nTargets:\n"} /^[a-zA-Z_-]+:.*?##/ { printf " \033[36m%-15s\033[0m %s\n", $$1, $$2 }' $(MAKEFILE_LIST)
