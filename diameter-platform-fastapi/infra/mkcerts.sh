#!/usr/bin/env bash
set -e
# Demo CA + server cert creation (NOT for production)
mkdir -p certs && cd certs
openssl genrsa -out ca.key 4096
openssl req -x509 -new -nodes -key ca.key -sha256 -days 3650 -subj "/CN=diameter-ca" -out ca.crt
openssl genrsa -out server.key 2048
openssl req -new -key server.key -subj "/CN=diameter-local.example" -out server.csr
openssl x509 -req -in server.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out server.crt -days 365 -sha256
echo "Created certs: $(pwd)/server.crt, server.key, ca.crt"
