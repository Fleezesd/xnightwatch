include scripts/make-rules/common.mk # make sure include common.mk at the first include line
include scripts/make-rules/all.mk

.PHONY: protoc
protoc: ## Generate api proto files.
	$(MAKE) gen.protoc