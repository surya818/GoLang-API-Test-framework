.DEFAULT_GOAL := test

# Ensure golang is available
ifeq (, $(shell which go 2> /dev/null))
$(error "'go' is not installed or available in PATH")
endif

APP_NAME := candidate-take-home-exercise-sdet
APP_DIR := $(patsubst %/,%,$(dir $(abspath $(lastword $(MAKEFILE_LIST)))))
APP_WORKDIR := $(shell pwd)

include $(APP_DIR)/mk/build.mk
ifeq ($(APP_DOCKER_BUILD),)
include $(APP_DIR)/mk/common.mk
include $(APP_DIR)/mk/docker.mk
include $(APP_DIR)/mk/tools.mk
endif
