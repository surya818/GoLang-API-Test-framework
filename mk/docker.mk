# --------------------------------------------------
# Docker tooling
# --------------------------------------------------

# Ensure docker is available
ifeq (, $(shell which docker 2> /dev/null))
$(error "'docker' is not installed or available in PATH")
endif

.PHONY: docker-build
docker-build:
	@docker build \
		-f "$(APP_DIR)/Dockerfile" \
		--build-arg APP_COMMIT=$(APP_COMMIT) \
		-t candidate-take-home-exercise-sdet \
		.

.PHONY: docker-run
docker-run: docker-build
	@docker run \
		-p 18080:18080 \
		--rm \
		--name candidate-app \
		-v "$(APP_DIR)/config.yml:/app/config.yml" \
		candidate-take-home-exercise-sdet
