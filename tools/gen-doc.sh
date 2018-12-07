#!/bin/bash

## 1. Please install swag first
## go get -u github.com/swaggo/swag/cmd/swag
##
## 2. Then, add ${GOPATH}/bin to PATH env

cd "${GOPATH}/src/yuudidi.com"
swag init -g cmd/server/main.go
