#!/bin/bash


go version
go env -w GO111MODULE=on
go env -w GOPROXY=https://goproxy.cn,direct
go env | grep GOPROXY
go mod tidy

GOOS=linux GOARCH=amd64 go build -o kano ./cmd/server/main.go

