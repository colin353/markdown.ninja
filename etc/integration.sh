#!/bin/bash

# If we're not running this integration on circleci, the $CIRCLE_SHA1
# won't be defined. So we'll define it to be a general local test tag.
TEST_DOCKER_TAG=${CIRCLE_SHA1:-test}

echo "USING TAG:"
echo $TEST_DOCKER_TAG

# Run the docker container so we can do our tests.
IMAGEPID=`docker run -d -e "APPCONFIG_MODE=test" --net=host colinmerkel/portfolio:$TEST_DOCKER_TAG`

# Presumably we run our integration tests here.
echo "Running jest..."
cd app/ && npm run test

# Stop the server. This is necessary because if the server crashed,
# this will return a nonzero exit code and fail the test.
docker stop $IMAGEPID
