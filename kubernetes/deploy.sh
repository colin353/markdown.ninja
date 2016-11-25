#!/bin/bash

# Exit on any error
set -e

export GOOGLE_APPLICATION_CREDENTIALS=~/account-auth.json

sudo chmod 777 /home/ubuntu/.kube/config

# Update Kubernetes replicationController
envsubst < kubernetes/portfolio.rc.yaml | kubectl create -f -

# Roll over Kubernetes pods
kubectl rolling-update portfolio
