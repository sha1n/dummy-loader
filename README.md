[![Build Status](https://travis-ci.com/sha1n/dummy-loader.svg?branch=master)](https://travis-ci.com/sha1n/dummy-loader) [![Go Report Card](https://goreportcard.com/badge/sha1n/dummy-loader)](https://goreportcard.com/report/sha1n/dummy-loader) [![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)


# Dummy Loader
Dummy loader is a very simple server written in Go that exposes an HTTP API for generating CPU and memory load. It was created for experimentation with resource management configuration in docker/kubernetes environments.


## Run The Dockerized Server
To run the server in a docker container, all you have to do is run the following command: 
```bash
docker run -p 8080:8080 sha1n/dummy-loader
```

## Generating CPU / Memory Load
1. Start CPU load by running `curl -v -XPOST http://<host>:<port>/api/cpu-load?time-sec=30[&cores=2]`
    1. `time-sec` - mandatory, the number of seconds to run the load
    2. `cores` - optional, how many cores to load 
2. Allocate memory on the server by running `curl -v -XPOST http://<host>:<port>/api/mem-footprint?amount-mb=1000`
    1. `amount-mb` - the amount of memory to allocate in MB.
    
*For convenience, these two endpoints support both HTTP GET and POST methods*     
   
