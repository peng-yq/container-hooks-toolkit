include $(CURDIR)/versions.mk

PREFIX := $(CURDIR)/bin
MODULE := container-hooks-toolkit

CMDS := $(patsubst ./cmd/%/,%,$(sort $(dir $(wildcard ./cmd/*/))))
CMD_TARGETS := $(patsubst %,bin-%, $(CMDS))

TESTS := $(notdir $(basename $(shell find ./test/* -type d)))
TEST_TARGETS := $(patsubst %,bin-%, $(TESTS))

ifeq ($(VERSION),)
CLI_VERSION = $(LIB_VERSION)$(if $(LIB_TAG),-$(LIB_TAG))
else
CLI_VERSION = $(VERSION)
endif
CLI_VERSION_PACKAGE = container-hooks-toolkit/internal/info

GOOS ?= linux

all: cmd test

ifneq ($(PREFIX),)
bin-%: COMMAND_BUILD_OPTIONS = -o $(PREFIX)/$(*)
endif

test: $(TEST_TARGETS)
$(TEST_TARGETS): bin-%:
	GOOS=$(GOOS) go build -ldflags "-extldflags=-Wl,-z,lazy -s -w -X $(CLI_VERSION_PACKAGE).gitCommit=$(GIT_COMMIT) -X $(CLI_VERSION_PACKAGE).version=$(CLI_VERSION)" $(COMMAND_BUILD_OPTIONS) $(MODULE)/test/$(*)

cmd: $(CMD_TARGETS)
$(CMD_TARGETS): bin-%:
	GOOS=$(GOOS) go build -ldflags "-extldflags=-Wl,-z,lazy -s -w -X $(CLI_VERSION_PACKAGE).gitCommit=$(GIT_COMMIT) -X $(CLI_VERSION_PACKAGE).version=$(CLI_VERSION)" $(COMMAND_BUILD_OPTIONS) $(MODULE)/cmd/$(*)

fmt:
	go fmt ./...

clean:
	rm -rf ./bin
