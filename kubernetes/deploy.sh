#!/bin/bash

# Exit on any error
set -e

sudo chmod 777 /home/ubuntu/.kube/config

# Update Kubernetes replicationController
envsubst < kubernetes/portfolio.rc.yaml | kubectl create -f -

# Roll over Kubernetes pods
kubectl rolling-update portfolio
