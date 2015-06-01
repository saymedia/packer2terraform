#!/bin/bash


# Cleanup
rm -f crash.log
rm -f p2t


# Testing
go test ./...
if [[ $? != 0 ]]; then
    exit 1
fi


# Linting
go vet ./...
if [[ $? != 0 ]]; then
    exit 1
fi

golint ./...
if [[ $? != 0 ]]; then
    exit 1
fi

gocyclo -over 15 .
if [[ $? != 0 ]]; then
    exit 1
fi


# Formatting
gofmt -s -d -l */*.go
if [[ $? != 0 ]]; then
    exit 1
fi


# Build
go build -o p2t
if [[ $? != 0 ]]; then
    exit 1
fi
