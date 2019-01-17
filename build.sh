#!/bin/bash

rm dist/*

proj_dir=`go env GOPATH`/src/yuudidi.com
cd $proj_dir
git pull origin dev

export GOOS=linux
export GOARCH=amd64

go build -o dist/webportal cmd/server/webportal/main.go
go build -o dist/app cmd/server/app/main.go
go build -o dist/h5backend cmd/server/h5backend/main.go
go build -o dist/background cmd/server/background/main.go
go build -o dist/websocket cmd/server/websocket/main.go

cp configs/config.yml dist/config.yml

zip -r dist.zip dist/*
