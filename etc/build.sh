#!/bin/bash

set -e

# Run the pre-build tasks.
cd app/ && node build.js

# Build the javascript minified binary.
