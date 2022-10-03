MAKEFILE_DIR := $(dir $(lastword $(MAKEFILE_LIST)))

test:
	@echo $(MAKEFILE_DIR)
test-run:

help:
	@echo ""