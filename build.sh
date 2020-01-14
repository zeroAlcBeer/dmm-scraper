#!/usr/bin/env bash

go build -ldflags '-w -s' .
./better-av-tool.exe -output "G:\Excited\failed\output" -path "G:\Excited\test" -proxy "socks5://127.0.0.1:7891"
#cp -f better-av-tool.exe "G:\Excited\JAV_output\better-av-tool.exe"