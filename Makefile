# Makefile
# https://www.alexedwards.net/blog/a-time-saving-makefile-for-your-go-projects

# Variables
BUILD_DIR = ./bin
CMD_DIR = ./cmd
APP_DIR = linter
APP_NAME = jl

# Set default goal (when you run make without specifying a goal)
.DEFAULT_GOAL := help

# ==================================================================================== #
# HELPERS
# ==================================================================================== #

## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

.PHONY: confirm
confirm:
	# Throws error if not confirmed
	@echo -n 'Are you sure? [y/N] ' && read ans && [ $${ans:-N} = y ]

.PHONY: no-dirty
no-dirty:
	# Check for uncommitted changes
	git diff --exit-code


# ==================================================================================== #
# QUALITY CONTROL
# ==================================================================================== #

## tidy: format code and tidy modfile
.PHONY: tidy
tidy:
	go fmt ./...
	go mod tidy -v

## audit: run quality control checks
.PHONY: audit
audit:
	go mod verify
	go vet ./...
	go run honnef.co/go/tools/cmd/staticcheck@latest -checks=all,-ST1000,-U1000 ./...
	go run golang.org/x/vuln/cmd/govulncheck@latest ./...
	go test -race -buildvcs -vet=off ./...


# ==================================================================================== #
# DEPENDENCIES
# ==================================================================================== #

## deps: download dependencies
.PHONY: deps
deps:
	go mod download


# ==================================================================================== #
# TESTING
# ==================================================================================== #

## test: run all tests
.PHONY: test
test:
	go test -v -race -buildvcs ./...

## test/cover: run all tests and display coverage
.PHONY: test/cover
test/cover:
	go test -v -race -buildvcs -coverprofile=/tmp/coverage.out ./...
	go tool cover -html=/tmp/coverage.out

## lint: run linter
.PHONY: lint
lint:
	golangci-lint run ./...


# ==================================================================================== #
# COMPILING
# ==================================================================================== #

## build: build all executables
.PHONY: build
build:
	@echo "Building executables..."
	mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(APP_NAME) $(CMD_DIR)/$(APP_DIR)

## build/scratch: build the scratch file located in ./scratch/scratch.go
.PHONY: build/scratch
build/scratch:
	@echo "Building scratch executable..."
	mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/scratch ./scratch


# ==================================================================================== #
# RUNNING
# ==================================================================================== #

## run: build and run executable
.PHONY: run
run: build
	@echo "Executing application..."
	@$(BUILD_DIR)/$(APP_NAME)

## run/scratch: build and run scratch executable
.PHONY: run/scratch
run/scratch: build/scratch
	@echo "Executing scratch..."
	@$(BUILD_DIR)/scratch


# ==================================================================================== #
# LIVE RELOADING
# ==================================================================================== #

## run/live: run parser application w/ live reloading
.PHONY: run/live
run/live:
	go run github.com/cosmtrek/air@latest \
		--build.cmd "make build" \
		--build.bin "$(BUILD_DIR)/$(APP_NAME)" \
		--build.delay "100" \
		--build.exclude_dir "" \
		--build.include_ext "go, tpl, tmpl, html, css, scss, js, ts, sql, jpeg, jpg, gif, png, bmp, svg, webp, ico" \
		--misc.clean_on_exit "true"


# ==================================================================================== #
# CLEANING
# ==================================================================================== #

## clean: clean up all build artifacts
.PHONY: clean
clean:
	@echo "Cleaning up..."
	go clean
	@rm -rf $(BUILD_DIR)/*


# ==================================================================================== #
# GIT
# ==================================================================================== #

## push: push changes to the remote Git repository, after running quality control checks
.PHONY: push
push: tidy audit no-dirty
	git push
