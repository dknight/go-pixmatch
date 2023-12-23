#!/bin/bash

BINDIR="./bin"
CMDNAME="pixmatch"
OS=$( uname )

mkdir -pv "$BINDIR"

export CGO_ENABLED='0'

echo "Compiling..."

GOOS=linux GOARCH=amd64 go build -o "$BINDIR/$CMDNAME-amd64-linux"
GOOS=linux GOARCH=arm64 go build -o "$BINDIR/$CMDNAME-arm64-linux"
GOOS=darwin GOARCH=amd64 go build -o "$BINDIR/$CMDNAME-amd64-darwin"
GOOS=darwin GOARCH=arm64 go build -o "$BINDIR/$CMDNAME-arm64-darwin"
GOOS=windows GOARCH=amd64 go build -o "$BINDIR/$CMDNAME.exe"

if [ "$OS" = "Linux" ]; then
    chmod +x "$BINDIR/$CMDNAME-amd64-linux"
fi

if [ "$OS" = "Darwin" ]; then
    chmod +x "$BINDIR/$CMDNAME-amd64-darwin"
    chmod +x "$BINDIR/$CMDNAME-arm64-darwin"
fi
