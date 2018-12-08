#!/bin/sh
#go run $(go env GOROOT)/src/crypto/tls/generate_cert.go --host=localhost
go run $(go env GOROOT)/src/crypto/tls/generate_cert.go --host=example.com
