#!/usr/bin/env bash
GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -trimpath -ldflags '-w -s'  -o bin/better-av-tool.exe .