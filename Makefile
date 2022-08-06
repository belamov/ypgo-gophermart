#!/usr/bin/make
# Makefile readme (ru): <http://linux.yaroslavl.ru/docs/prog/gnu_make_3-79_russian_manual.html>
# Makefile readme (en): <https://www.gnu.org/software/make/manual/html_node/index.html#SEC_Contents>

SHELL = /bin/sh

app_container_name := gophermart
docker_bin := $(shell command -v docker 2> /dev/null)
docker_compose_bin := $(shell command -v docker-compose 2> /dev/null)
docker_compose_yml := docker/docker-compose.yml
user_id := $(shell id -u)

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
	$(docker_compose_bin) --file "$(docker_compose_yml)" build

up: build ## Run app
	$(docker_compose_bin) --file "$(docker_compose_yml)" up

mock: ## Generate mocks
	$(docker_compose_bin) --file "$(docker_compose_yml)" run --rm $(app_container_name) bash -c "\
 		mockgen -destination=internal/gophermart/mocks/auth.go -package=mocks github.com/belamov/ypgo-gophermart/internal/gophermart/services Authenticator && \
		mockgen -destination=internal/gophermart/mocks/users.go -package=mocks github.com/belamov/ypgo-gophermart/internal/gophermart/storage Users \
		"

lint:
	$(docker_bin) run --rm -v $(shell pwd):/app -w /app golangci/golangci-lint:v1.46.2 golangci-lint run

gofumpt:
	$(docker_compose_bin) --file "$(docker_compose_yml)" run --rm $(app_container_name) gofumpt -l -w .

test: ## Execute tests
	$(docker_compose_bin) --file "$(docker_compose_yml)" run --rm $(app_container_name) go test -v ./...

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