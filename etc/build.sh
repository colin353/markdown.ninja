#!/bin/bash

set -e

cd ./app/

# Create the HTML renders of the markdown files.
node etc/build.js

# Prerender prescribed routes to HTML.
./node_modules/.bin/babel-node etc/prerender.js

# Create the minified JS
echo "Creating minified production JS..."
export NODE_ENV=production
trap 'kill %1' SIGINT
node_modules/.bin/webpack --optimize-minimize --optimize-dedupe --progress --colors -p --config webpack.config.prod.js

# Compile the docker-compatible Go binary
echo "Compiling special go binary..."
cd ..
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo .
