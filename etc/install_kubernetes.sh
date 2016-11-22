#!/bin/bash

set -e

if [ -d "~/kubernetes" ]; then
  # If the directory exists, there is no need to refresh the cache.
  exit 0
fi

curl -o ~/kubernetes/kubectl https://storage.googleapis.com/kubernetes-release/release/v1.4.3/bin/linux/amd64/kubectl
chmod +x ~/kubernetes/kubectl
