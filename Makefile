GO := $(shell which go)
IMAGE_TAG := 0.1.0
main_package_path = ./cmd/authProxy
binary_name = authProxy

# ==================================================================================== #
# HELPERS
# ==================================================================================== #

.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

# ==================================================================================== #
# DEVELOPMENT
# ==================================================================================== #

.PHONY: build

## build: build the application
build:
    # Include additional build steps, like TypeScript, SCSS or Tailwind compilation here...
	CGO_ENABLED=1 $(GO) build -o=bin/${binary_name} ${main_package_path}

## run: run the  application
.PHONY: run
run: build
	bin/${binary_name}

## run/live: run the application with reloading on file changes
.PHONY: dev
dev:
	bin/air \
	--build.cmd "make build" --build.bin "./bin/${binary_name}" --build.delay "100" \
	--build.exclude_dir "" \
	--build.include_ext "go, tpl, tmpl, html, css, scss, js, ts, sql, jpeg, jpg, gif, png, bmp, svg, webp, ico" \
	--misc.clean_on_exit "true"

## up.dev: run docker compose
.PHONY: up.dev
up.dev:
	TAG=$(IMAGE_TAG) ENV_FILE=../.env docker compose -f deployments/docker-compose.yml up -d

## push.dev: push docker image
.PHONY: push.dev
push.dev:
	TAG=$(IMAGE_TAG) ENV_FILE=../.env docker compose -f deployments/docker-compose.yml build --no-cache
	TAG=$(IMAGE_TAG) ENV_FILE=../.env docker compose -f deployments/docker-compose.yml push

## down.dev: stop docker compose
.PHONY: down.dev
down.dev:
	TAG=$(IMAGE_TAG) ENV_FILE=../.env docker compose -f deployments/docker-compose.yml down

## tidy: resync modules
.PHONY: tidy
tidy:
	$(GO) mod tidy

