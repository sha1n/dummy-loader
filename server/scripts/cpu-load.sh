#!/usr/bin/env bash

curl -v -XPOST "http://localhost:8080/api/cpu-load?time-sec=$1&cores=$2"