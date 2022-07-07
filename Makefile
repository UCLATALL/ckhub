.DEFAULT_GOAL = help

BUILD_DIR := $(if $(BUILD_DIR),$(BUILD_DIR:/=),target)

GIT_SHA    := $(shell git rev-parse --short HEAD 2>/dev/null)
GIT_TAG    := $(shell git tag --points-at HEAD 2>/dev/null)

VERSION ?= $(if $(GIT_TAG),$(GIT_TAG:v%=%),unknown$(if $(GIT_SHA), on $(GIT_SHA),))

BIN_DIR  := $(BUILD_DIR)/bin
GOFLAGS  ?= -trimpath
LDFLAGS  ?= -X 'main.version=$(VERSION)' -extldflags -static
PLATFORM ?= linux/amd64 linux/arm64
TARGET   ?= $(patsubst cmd/%/main.go,%,$(wildcard cmd/*/main.go))
IMAGE    ?= us-central1-docker.pkg.dev/ckhub-proto1/ckhub/play

OPENAPI_DIR := $(BUILD_DIR)/openapi
REPORTS_DIR := $(BUILD_DIR)/reports

## build: Build application binaries.
builds := $(foreach platform,$(PLATFORM),$(TARGET:%=%/$(platform)))

.PHONY: build
build: $(builds)

$(BIN_DIR):
	mkdir -p $@

.PHONY: $(builds)
$(builds): | $(BIN_DIR)
	$(eval $@:source := ./cmd/$(word 1,$(subst /, ,$@)))
	$(eval $@:goos   := $(word 2,$(subst /, ,$@)))
	$(eval $@:goarch := $(word 3,$(subst /, ,$@)))
	$(eval $@:goarm  := $(word 4,$(subst /, ,$@)))
	$(eval $@:goenv  := $(if $(goos),GOOS=$(goos) ,))
	$(eval $@:goenv  := $(goenv)$(if $(goarch),GOARCH=$(goarch) ,))
	$(eval $@:goenv  := $(goenv)$(if $(goarm),GOARM=$(goarm) ,))
	$(eval $@:output := $(BIN_DIR)/$(subst /,-,$@)$(if $(goos:windows=),,.exe))
	$(eval $@:flags  := $(GOFLAGS)$(if $(LDFLAGS), -ldflags "$(LDFLAGS)",))
	$(goenv)go build $(flags) -o $(output) $(source)

## clean: Remove created resources.
.PHONY: clean
clean:
	rm -rf $(BIN_DIR) $(REPORTS_DIR) $(OPENAPI_DIR)

## docker: Create docker images.
images := $(foreach platform,$(PLATFORM),docker/$(platform))

.PHONY: docker
docker: $(images)

.PHONY: $(images)
$(images): | build
	$(eval $@:platform := $(@:docker/%=%))
	$(eval $@:branch   := $(subst /,-,$(if $(GIT_BRANCH),$(GIT_BRANCH),unknown)))
	$(eval $@:image    := $(IMAGE):$(branch)-$(subst /,-,$(platform)))
	docker build -f .docker/Dockerfile -t $(image) --platform $(platform) $(BUILD_DIR)

## help: Display available targets.
.PHONY: help
help: $(MAKEFILE_LIST)
	@echo "Usage: make [target]"
	@echo
	@echo "Targets:"
	@sed -En 's/^## *([^:]+): *(.*)$$/\1\t\2/p' $< | expand -t 18

## lint: Run static analysis checks.
.PHONY: lint
lint:
	golangci-lint run

## openapi: Generate openapi specification.
.PHONY: openapi
openapi: | swagger
	curl -so $(OPENAPI_DIR)/play.yml -XPOST -d @target/openapi/play.json \
	-H 'accept: application/yaml' -H 'Content-Type: application/json' \
	https://converter.swagger.io/api/convert

swagger: | $(OPENAPI_DIR)
	swagger generate spec -m --output $(OPENAPI_DIR)/play.json

$(OPENAPI_DIR):
	mkdir -p $@

## test: Run tests and generate quality reports.
tests := $(PLATFORM:%=test/go/%)

.PHONY: test
test: $(tests)

.PHONY: $(tests)
$(tests): | $(REPORTS_DIR)
	$(eval $@:goos   := $(word 1,$(subst /, ,$(@:test/go/%=%))))
	$(eval $@:goarch := $(word 2,$(subst /, ,$(@:test/go/%=%))))
	$(eval $@:goarm  := $(word 3,$(subst /, ,$(@:test/go/%=%))))
	$(eval $@:goenv  := $(if $(goos),GOOS=$(goos) ,))
	$(eval $@:goenv  := $(goenv)$(if $(goarch),GOARCH=$(goarch) ,))
	$(eval $@:goenv  := $(goenv)$(if $(goarm),GOARM=$(goarm) ,))
	$(eval $@:suffix := $(subst /,-,$(@:test/go/%=%)))
	$(goenv)gotestsum --junitfile $(REPORTS_DIR)/tests-$(suffix).xml -f standard-quiet -- \
	-coverpkg ./... -covermode atomic -coverprofile $(REPORTS_DIR)/cover-$(suffix).out \
	./...
	coverage -i "$(wildcard cmd/*/*.go)" $(REPORTS_DIR)/cover-$(suffix).out

$(REPORTS_DIR):
	mkdir -p $@

## version: Display current version.
.PHONY: version
version:
	@echo 'version $(VERSION)'
