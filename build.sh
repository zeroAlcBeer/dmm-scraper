#!/usr/bin/env bash

go build -ldflags '-w -s' .
./upx.exe better-av-tool.exe