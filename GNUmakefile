TEST?=$$(go list ./...)
GOFMT_FILES?=$$(find . -name '*.go')
WEBSITE_REPO=github.com/hashicorp/terraform-website

PKG_NAME=nhncloud
VERSION_STR=1.0.7

OS_IMPLEMENTATION=$(shell uname -s | tr '[:upper:]' '[:lower:]')
HW_PLATFORM=$(shell uname -m)
OS_HW_STR=$(OS_IMPLEMENTATION)_$(HW_PLATFORM)

GO_BIN_PATH=$(GOPATH)/bin/terraform-provider-$(PKG_NAME)
TF_BIN_PATH=~/.terraform.d/plugins/terraform.local/local/$(PKG_NAME)/$(VERSION_STR)/$(OS_HW_STR)/terraform-provider-$(PKG_NAME)_v$(VERSION_STR)

platforms := darwin/amd64 darwin/arm64 \
			 linux/386 linux/amd64 linux/arm linux/arm64 \
			 windows/386 windows/amd64 \
			freebsd/386 freebsd/amd64

all: fmtcheck
	@if [ ! -d "bin" ]; then\
		mkdir bin;\
	fi
	$(foreach platform, $(platforms), GOOS=$(word 1, $(subst /, ,$(platform))) GOARCH=$(word 2, $(subst /, ,$(platform))) go build -o bin/$(subst /,_,$(platform))/terraform-provider-$(PKG_NAME)_v$(VERSION_STR);)

default: build

build: fmtcheck
	go install

local: build
	@echo "\n\033[4m\033[96mInstalling locally...\033[0m"
	@echo "\033[97m | Operating system implementation : $(OS_IMPLEMENTATION)\033[0m"
	@echo "\033[97m | Haredware platform              : $(HW_PLATFORM)\033[0m"
	@echo "\033[97m | Package name                    : $(PKG_NAME)\033[0m"
	@echo "\033[97m | Version string                  : $(VERSION_STR)\033[0m\n"
	@cp -v $(GO_BIN_PATH) $(TF_BIN_PATH)

test: fmtcheck
	go test -i $(TEST) || exit 1
	echo $(TEST) | \
		xargs -t -n4 go test $(TESTARGS) -timeout=30s -parallel=4

testacc: fmtcheck
	TF_ACC=1 TF_ACC_TERRAFORM_VERSION=1.2.9 go test $(TEST) -v $(TESTARGS) -timeout 120m

vet:
	@echo "go vet ."
	@go vet $$(go list ./...) ; if [ $$? -eq 1 ]; then \
		echo ""; \
		echo "Vet found suspicious constructs. Please check the reported constructs"; \
		echo "and fix them if necessary before submitting the code for review."; \
		exit 1; \
	fi

fmt:
	gofmt -w $(GOFMT_FILES)

fmtcheck:
	@sh -c "'$(CURDIR)/scripts/gofmtcheck.sh'"

errcheck:
	@sh -c "'$(CURDIR)/scripts/errcheck.sh'"

test-compile:
	@if [ "$(TEST)" = "./..." ]; then \
		echo "ERROR: Set TEST to a specific package. For example,"; \
		echo "  make test-compile TEST=./$(PKG_NAME)"; \
		exit 1; \
	fi
	go test -c $(TEST) $(TESTARGS)

website:
ifeq (,$(wildcard $(GOPATH)/src/$(WEBSITE_REPO)))
	echo "$(WEBSITE_REPO) not found in your GOPATH (necessary for layouts and assets), get-ting..."
	git clone https://$(WEBSITE_REPO) $(GOPATH)/src/$(WEBSITE_REPO)
endif
	@$(MAKE) -C $(GOPATH)/src/$(WEBSITE_REPO) website-provider PROVIDER_PATH=$(shell pwd) PROVIDER_NAME=$(PKG_NAME)

website-test:
ifeq (,$(wildcard $(GOPATH)/src/$(WEBSITE_REPO)))
	echo "$(WEBSITE_REPO) not found in your GOPATH (necessary for layouts and assets), get-ting..."
	git clone https://$(WEBSITE_REPO) $(GOPATH)/src/$(WEBSITE_REPO)
endif
	@$(MAKE) -C $(GOPATH)/src/$(WEBSITE_REPO) website-provider-test PROVIDER_PATH=$(shell pwd) PROVIDER_NAME=$(PKG_NAME)

# Linux 전용 빌드 타겟
linux: fmtcheck
	@echo "\n\033[4m\033[96mBuilding for Linux platforms...\033[0m"
	@if [ ! -d "dist" ]; then\
		mkdir -p dist;\
	fi
	@echo "\033[97m | Building for linux/amd64...\033[0m"
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-s -w" -trimpath -o dist/linux_amd64/terraform-provider-$(PKG_NAME)_v$(VERSION_STR)
	@echo "\033[97m | Building for linux/arm64...\033[0m"
	GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -ldflags="-s -w" -trimpath -o dist/linux_arm64/terraform-provider-$(PKG_NAME)_v$(VERSION_STR)
	@echo "\033[92m✓ Linux builds completed successfully!\033[0m"

# Linux 빌드 스크립트 실행
linux-script:
	@./build-linux.sh

.PHONY: build test testacc vet fmt fmtcheck errcheck test-compile website website-test local linux linux-script

