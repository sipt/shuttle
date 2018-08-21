#!/usr/bin/env bash

if [ ! -d "_output/" ];then
mkdir _output
else
rm -rf _output/*
fi

cp -rf view _output/
cp GeoLite2-Country.mmdb _output/
if [ -d ".conf/" ];then
cp .conf/sipt.yaml _output/shuttle.yaml
else
cp example.yaml _output/shuttle.yaml
fi
go build -o _output/shuttle cmd/main.go

echo "nohup ./shuttle >> shuttle.log 2>&1 &" > _output/start.sh
chmod a+x _output/start.sh


