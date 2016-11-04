#!/bin/bash

trap 'kill %1' SIGINT
webpack --progress --colors --watch
