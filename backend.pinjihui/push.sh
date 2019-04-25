#!/bin/sh
go generate ./schema &&
GOOS=linux GOARCH=amd64 go build -ldflags "-w -s" -o backend.pinjihui &&
git add . &&
git commit -m "$(date +%Y-%m-%d' '%H:%M:%S)" &&
git push