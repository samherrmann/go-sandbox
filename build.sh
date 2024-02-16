#!/bin/bash
mkdir -p dist
go build -o dist -ldflags "-w -s" .
