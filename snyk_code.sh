#!/bin/bash

sudo docker run --rm \
    -e SNYK_TOKEN=${SNYK_TOKEN} \
    -v "${PWD}":/app \
    snyk/snyk:golang-1.20 snyk test --severity-threshold=high

sudo docker run --rm \
    -e SNYK_TOKEN=${SNYK_TOKEN} \
    -v "${PWD}/frontend":/app \
    snyk/snyk:node-18 snyk test --severity-threshold=high
