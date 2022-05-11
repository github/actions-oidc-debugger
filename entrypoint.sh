#!/bin/sh -l

cd /

go run cmd/oidc-debug.go -audience $1
