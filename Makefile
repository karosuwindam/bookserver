MAKEFILE_DIR := $(dir $(lastword $(MAKEFILE_LIST)))

test:
	@echo $(MAKEFILE_DIR)
test-run:

create:
	@echo create
build: create
	@echo build
help:
	@echo ""