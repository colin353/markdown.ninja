#!/bin/bash

# Exit on any error
set -e

# Tell Google where to find the account credentials.
export GOOGLE_APPLICATION_CREDENTIALS=~/account-auth.json

# For some reaosn we need to adjust the permissions of the
# .kube folder, I think maybe because we ran some previous
# commands as root (?)
sudo chmod -R 777 /home/ubuntu/.kube

# Roll over Kubernetes pods
kubectl rolling-update portfolio --image=colinmerkel/portfolio:${CIRCLE_SHA1}
