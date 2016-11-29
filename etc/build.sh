#!/bin/bash

set -e

cd ./app/

# Create the HTML renders of the markdown files.
echo "--- Loading CSS ---"
node etc/indexcss.js

# Create the HTML renders of the markdown files.
echo "--- Render markdown ---"
node etc/build.js

echo "--- Prerender react components ---"
# Prerender prescribed routes to HTML.
./node_modules/.bin/babel-node etc/prerender.js

# Create the minified JS
echo "--- Minify production JS ---"
export NODE_ENV=production
trap 'kill %1' SIGINT
node_modules/.bin/webpack --optimize-minimize --bail --optimize-dedupe --progress --colors -p --config webpack.config.prod.js

# Compile the docker-compatible Go binary. It's statically
# linked so that it can run inside of a minimalist docker
# container.
echo "--- Compile special go binary for docker environment ---"
cd ..
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo .

echo "--- Build complete. ---"
