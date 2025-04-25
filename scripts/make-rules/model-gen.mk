include scripts/make-rules/common.mk


.PHONY: compose-up-model-gen
compose-up-model-gen:
	docker-compose -f build/docker/gen-gorm-model/docker-compose.yaml up -d

.PHONY: compose-down-model-gen
compose-down-model-gen:
	docker-compose -f build/docker/gen-gorm-model/docker-compose.yaml down