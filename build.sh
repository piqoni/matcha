#!/bin/bash

GOOS=linux GOARCH=amd64 go build -o bin/matcha-amd64-linux .
GOOS=darwin GOARCH=amd64 go build -o bin/matcha-amd64-darwin .
GOOS=windows GOARCH=amd64 go build -o bin/matcha-amd64.exe .