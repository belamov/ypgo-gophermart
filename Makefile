#!/usr/bin/make
# Makefile readme (ru): <http://linux.yaroslavl.ru/docs/prog/gnu_make_3-79_russian_manual.html>
# Makefile readme (en): <https://www.gnu.org/software/make/manual/html_node/index.html#SEC_Contents>

SHELL = /bin/sh

app_container_name := gophermart
docker_bin := $(shell command -v docker 2> /dev/null)
docker_compose_bin := $(shell command -v docker-compose 2> /dev/null)
docker_compose_yml := docker/docker-compose.yml
user_id := $(shell id -u)
project_name := ypgo_gophermart

.PHONY : help pull build push login test clean \
         app-pull app app-push\
         sources-pull sources sources-push\
         nginx-pull nginx nginx-push\
         up down restart shell install
.DEFAULT_GOAL := help

# --- [ Development tasks ] -------------------------------------------------------------------------------------------
help:  ## Display this help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n\nTargets:\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-10s\033[0m %s\n", $$1, $$2 }' $(MAKEFILE_LIST)

build: ## Build containers
	$(docker_compose_bin) --file "$(docker_compose_yml)" build --parallel

up: build ## Run app
	$(docker_compose_bin) --file "$(docker_compose_yml)" up --remove-orphans

mocks_delete:
	$(docker_compose_bin) --file "$(docker_compose_yml)" run  --rm $(app_container_name) bash -c "rm -r -f internal/gophermart/mocks/ internal/accrual/mocks/"

mocks_regenerate: mocks_delete mocks_generate

mocks_generate: ## Generate mocks
	$(docker_compose_bin) --file "$(docker_compose_yml)" run  --rm $(app_container_name) bash -c "\
 		mockgen -destination=internal/gophermart/mocks/auth.go -package=mocks github.com/belamov/ypgo-gophermart/internal/gophermart/services Auth && \
		mockgen -destination=internal/gophermart/mocks/orders_service.go -package=mocks github.com/belamov/ypgo-gophermart/internal/gophermart/services OrdersManagerInterface && \
		mockgen -destination=internal/gophermart/mocks/balance_service.go -package=mocks github.com/belamov/ypgo-gophermart/internal/gophermart/services BalanceProcessorInterface && \
		mockgen -destination=internal/gophermart/mocks/accrual_service.go -package=mocks github.com/belamov/ypgo-gophermart/internal/gophermart/services AccrualInfoProvider && \
		mockgen -destination=internal/gophermart/mocks/users_storage.go -package=mocks github.com/belamov/ypgo-gophermart/internal/gophermart/storage UsersStorage && \
		mockgen -destination=internal/gophermart/mocks/orders_storage.go -package=mocks github.com/belamov/ypgo-gophermart/internal/gophermart/storage OrdersStorage &&\
		mockgen -destination=internal/gophermart/mocks/balance_storage.go -package=mocks github.com/belamov/ypgo-gophermart/internal/gophermart/storage BalanceStorage  &&\
		\
		mockgen -destination=internal/accrual/mocks/orders_manager.go -package=mocks github.com/belamov/ypgo-gophermart/internal/accrual/services OrderManagementInterface  &&\
		mockgen -destination=internal/accrual/mocks/orders_storage.go -package=mocks github.com/belamov/ypgo-gophermart/internal/accrual/storage OrdersStorage  &&\
		mockgen -destination=internal/accrual/mocks/rewards_storage.go -package=mocks github.com/belamov/ypgo-gophermart/internal/accrual/storage RewardsStorage  \
		"

lint:
	$(docker_bin) run --rm -v $(shell pwd):/app -w /app golangci/golangci-lint:v1.46.2 golangci-lint run

gofumpt:
	$(docker_compose_bin) --file "$(docker_compose_yml)" run --rm $(app_container_name) gofumpt -l -w .

test: ## Execute tests
	$(docker_compose_bin) --file "$(docker_compose_yml)" run --rm $(app_container_name) go test -v ./internal/gophermart/...
	$(docker_compose_bin) --file "$(docker_compose_yml)" run --rm accrual go test -v ./internal/accrual/...

check: build gofumpt lint test  ## Run tests and code analysis

# Prompt to continue
prompt-continue:
	@while [ -z "$$CONTINUE" ]; do \
		read -r -p "Would you like to continue? [y]" CONTINUE; \
	done ; \
	if [ ! $$CONTINUE == "y" ]; then \
        echo "Exiting." ; \
        exit 1 ; \
    fi
