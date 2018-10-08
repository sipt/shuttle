#!/usr/bin/env bash

# ./build.sh mac|win|linux

if [ ! -d "_output/" ];then
mkdir _output
fi
if [ ! -d "_output/shuttle" ];then
mkdir _output/shuttle
else
rm -rf _output/shuttle/*
fi

cp -rf view _output/shuttle/
cp GeoLite2-Country.mmdb _output/shuttle/
#if [ -d ".conf/" ];then
#cp .conf/sipt.yaml _output/shuttle.yaml
#else
cp example.yaml _output/shuttle/shuttle.yaml
#fi
mkdir _output/shuttle/RespFiles

if [ "$1" == "mac" ];then
# mac os
GOOS=darwin GOARCH=amd64 go build -o _output/shuttle/shuttle cmd/main.go
GOOS=darwin GOARCH=amd64 go build -o _output/shuttle/upgrade scripts/upgrade.go
`echo "c
nohup ./shuttle >> shuttle.log 2>&1 &" > _output/shuttle/start.sh`
`chmod a+x _output/shuttle/start.sh`
elif [ "$1" == "win" ];then
# windows
GOOS=windows GOARCH=amd64 go build -o _output/shuttle/shuttle.exe cmd/main.go
GOOS=windows GOARCH=amd64 go build -o _output/shuttle/upgrade.exe scripts/upgrade.go
`echo "@echo off
if \"%1\" == \"h\" goto begin
mshta vbscript:createobject(\"wscript.shell\").run(\"%~nx0 h\",0)(window.close)&&exit
:begin
shuttle >> shuttle.log" > _output/shuttle/startup.bat`
elif [ "$1" == "linux" ];then
# linux
GOOS=linux GOARCH=amd64 go build -o _output/shuttle/shuttle cmd/main.go
GOOS=linux GOARCH=amd64 go build -o _output/shuttle/upgrade scripts/upgrade.go
`echo "#!/usr/bin/env bash
nohup ./shuttle >> shuttle.log 2>&1 &" > _output/shuttle/start.sh`
`chmod a+x _output/shuttle/start.sh`
fi
