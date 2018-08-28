#!/usr/bin/env bash

curl -v -XPOST "http://localhost:8080/api/mem-footprint?amount-mb=$1"