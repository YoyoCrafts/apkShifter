#!/bin/sh

rm -rf  apkShifter.zip

GOOS=linux GOARCH=amd64 go build -o apkShifter main.go

zip -r apkShifter.zip  ./config ./apkShifter ./run.sh


