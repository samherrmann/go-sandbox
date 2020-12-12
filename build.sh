#!/bin/bash
mkdir -p dist
go build -o dist ./app
go build -buildmode=plugin -o dist ./plugin
