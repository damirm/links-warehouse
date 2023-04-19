#!/usr/bin/env bash

set -ex

# Change dir to the root of repo
cd "$(dirname "${BASH_SOURCE[0]}")/../"

DOCKERFILE_PATH=$1
IMAGE_NAME=$2
IMAGE_TAG=$3
COMMAND=$4

if [[ "${COMMAND}" == "build" ]]; then
  echo "Building docker image: ${IMAGE_NAME}"
  docker build \
    -f "${DOCKERFILE_PATH}" \
    -t "${IMAGE_NAME}:${IMAGE_TAG}" .

  docker tag "${IMAGE_NAME}:${IMAGE_TAG}" "${IMAGE_NAME}:latest"
fi

if [[ "${COMMAND}" == "push" ]]; then
  docker push "${IMAGE_NAME}:${IMAGE_TAG}"
  docker push "${IMAGE_NAME}:latest"
fi
