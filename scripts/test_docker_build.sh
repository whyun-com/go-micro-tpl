#!/bin/bash

export GIT_BRANCH_FOR_MAKE=development
export AREA_TYPE=idc
export CI_COMMIT_TAG=0.0.0
export CI_PROJECT_NAME=go-micro-tpl
source server_env.sh

./docker_build.sh