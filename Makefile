include scripts/make-rules/common.mk # make sure include common.mk at the first include line
include scripts/make-rules/all.mk

.PHONY: gen
gen: ## Generate CI-related files. Generate all files by specifying `A=1`.
ifeq ($(ALL),1)
	$(MAKE) gen.all
else
	$(MAKE) gen.run
endif

.PHONY: install-tools
install-tools: ## Install CI-related tools. Install all tools by specifying `A=1`.
	$(MAKE) _install.ci
	if [[ "$(A)" == 1 ]]; then                                             \
		$(MAKE) _install.other ;                                            \
	fi

.PHONY: protoc
protoc: ## Generate api proto files.
	$(MAKE) gen.protoc
