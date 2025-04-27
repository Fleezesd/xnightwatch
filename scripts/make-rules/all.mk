# ==============================================================================
# Makefile used to aggregate all makefiles for easy management.
#

include scripts/make-rules/tools.mk
include scripts/make-rules/golang.mk
include scripts/make-rules/generate.mk
include scripts/make-rules/generate-k8s.mk
include scripts/make-rules/model-gen.mk