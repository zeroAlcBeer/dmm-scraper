#!/usr/bin/env bash

go build -ldflags '-w -s' .
./upx.exe better-av-tool.exe
#./better-av-tool-v0.7.1-win64.exe -output "G:\Excited\failed\output" -path "G:\Excited\test" -proxy "socks5://127.0.0.1:7891"
#cp -f better-av-tool-v0.7.1-win64.exe "G:\Excited\JAV_output\better-av-tool-v0.7.1-win64.exe"