#!/bin/bash

killall background
killall app
killall webportal
killall h5backend
killall websocket

logPath=${HOME}/log
mkdir -p ${logPath}

cd dist

nohup ./background 2>&1 > ${logPath}/background-$(date +%Y%m%d).log  &
export GIN_PORT=8081
nohup ./app 2>&1 > ${logPath}/app-$(date +%Y%m%d).log  &
export GIN_PORT=8082
nohup ./webportal 2>&1 > ${logPath}/webportal-$(date +%Y%m%d).log  &
export GIN_PORT=8083
nohup ./h5backend 2>&1 > ${logPath}/h5backend-$(date +%Y%m%d).log  &
export GIN_PORT=8085
nohup ./websocket 2>&1 > ${logPath}/websocket-$(date +%Y%m%d).log  &
