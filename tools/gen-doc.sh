#!/bin/bash

## 1. Please install swag first
## go get -u github.com/swaggo/swag/cmd/swag
##
## 2. Then, add ${GOPATH}/bin to PATH env

cd "${HOME}/go/src/YuuPay_core-service"
swag init -g cmd/server/main.go