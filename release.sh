#!/usr/bin/env bash

rm -rf _output/*.zip

echo "build for windows ..."
./build.sh win
cd _output/
zip -q -r shuttle_windows_amd64_$1.zip shuttle

cd ..
echo "build for linux ..."
./build.sh linux
cd _output/
zip -q -r shuttle_linux_amd64_$1.zip shuttle

cd ..
echo "build for mac ..."
./build.sh mac
cd _output/
zip -q -r shuttle_macos_amd64_$1.zip shuttle

echo "end ..."