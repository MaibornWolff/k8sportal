#!/usr/bin/env bash

make build
docker run -p 8080:8080 k8sportal
