#!/bin/bash

if [ -e "~/node_modules" ]; then
  mv ~/node_modules admin/
  echo "Restoring node_modules from existing cache."
else
  echo "Can't find node_modules cache."
fi
