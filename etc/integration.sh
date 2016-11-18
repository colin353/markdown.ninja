#!/bin/bash

IMAGEPID=`docker run -d --net=host colinmerkel/portfolio:$CIRCLE_SHA1`
# Presumably we run our integration tests here.
echo $IMAGEPID
# Stop the server. This is necessary because if the server crashed,
# this will return a nonzero exit code and fail the test.
docker stop $IMAGEPID
