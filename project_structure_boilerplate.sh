mkdir -p diameter-platform && cd diameter-platform

# create tree
mkdir -p go-diameter-server/cmd/server
mkdir -p go-diameter-server/internal/diameter
mkdir -p go-diameter-server/internal/api
mkdir -p go-diameter-server/internal/store
mkdir -p go-diameter-server/internal/config
mkdir -p go-diameter-server/internal/metrics
mkdir -p go-diameter-server/config

mkdir -p django-control/control
touch django-control/manage.py
mkdir -p infra
mkdir -p tests/integration

# top-level files
cat > .gitignore <<'EOF'
# Go
/bin/
/vendor/
/*.exe
*.log

# Python/Django
__pycache__/
*.pyc
db.sqlite3

# Docker
.env
docker-compose.override.yml

# IDEs
.vscode/
.idea/
EOF

# README
cat > README.md <<'EOF'
# Diameter Platform

Scaffold for a production-ready Diameter service: Go protocol server + Django control UI + infra (docker-compose, cert helper).

Usage:
  - See go-diameter-server/ and django-control/ for entry points.
  - Start local stack with docker-compose in infra/.
EOF

# go main.go
cat > go-diameter-server/cmd/server/main.go <<'EOF'
package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// TODO: load config, init logger, init stores, start diameter listener and http api
	log.Println("diameter server starting (scaffold)")

	// placeholder graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	<-ctx.Done()
	log.Println("shutting down")
}
EOF

# handler scaffold
cat > go-diameter-server/internal/diameter/handler.go <<'EOF'
package diameter

// Handler stubs for Diameter protocol handling.
// Implement CER/CEA, DWR/DWA, ACR/ACA, AAR/AAA here.

func InitHandlers() {
    // register handlers with go-diameter lib
}
EOF

cat > go-diameter-server/internal/diameter/wire.go <<'EOF'
package diameter

// wire helpers: AVP builders / parsers, command creation helpers
EOF

# API stubs
cat > go-diameter-server/internal/api/handlers.go <<'EOF'
package api

import "net/http"

// RegisterHandlers registers the management http handlers.
func RegisterHandlers(mux *http.ServeMux) {
    mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        w.Write([]byte("ok"))
    })
    mux.HandleFunc("/peers", func(w http.ResponseWriter, r *http.Request) {
        // TODO: list/add peers
        w.Write([]byte("peers endpoint"))
    })
}
EOF

# store stubs
cat > go-diameter-server/internal/store/redis.go <<'EOF'
package store

// Redis session store stub - implement session set/get and TTL
EOF

cat > go-diameter-server/internal/store/pg.go <<'EOF'
package store

// Postgres audit store stub - implement audit write and queries
EOF

# config loader
cat > go-diameter-server/internal/config/loader.go <<'EOF'
package config

// Config loader stub - parse YAML into config structs and validate
EOF

# metrics stub
cat > go-diameter-server/internal/metrics/metrics.go <<'EOF'
package metrics

// Prometheus metrics registration helpers
EOF

# example config
cat > go-diameter-server/config/example.yaml <<'EOF'
origin_host: diameter-local.example
origin_realm: example
listen:
  address: 0.0.0.0:3868
  tls: false
redis:
  addr: redis:6379
postgres:
  dsn: postgres://diam:diam@postgres:5432/diameter?sslmode=disable
peers: []
EOF

# Dockerfile for Go server
cat > go-diameter-server/Dockerfile <<'EOF'
FROM golang:1.20-alpine AS build
WORKDIR /src
COPY . .
RUN go build -o /bin/diameter ./cmd/server

FROM alpine:3.18
COPY --from=build /bin/diameter /usr/local/bin/diameter
COPY config /etc/diameter
EXPOSE 3868 8080
CMD ["/usr/local/bin/diameter"]
EOF

# django minimal
cat > django-control/manage.py <<'EOF'
#!/usr/bin/env python
import os
import sys
if __name__ == "__main__":
    os.environ.setdefault("DJANGO_SETTINGS_MODULE", "control.settings")
    from django.core.management import execute_from_command_line
    execute_from_command_line(sys.argv)
EOF

cat > django-control/control/settings.py <<'EOF'
SECRET_KEY = "replace-me"
DEBUG = True
ALLOWED_HOSTS = ["*"]
ROOT_URLCONF = "control.urls"
EOF

cat > django-control/control/urls.py <<'EOF'
from django.urls import path
from . import views
urlpatterns = [
    path("", views.index),
]
EOF

cat > django-control/control/views.py <<'EOF'
from django.http import HttpResponse
def index(request):
    return HttpResponse("Diameter Control UI (scaffold)")
EOF

cat > django-control/Dockerfile <<'EOF'
FROM python:3.11-slim
WORKDIR /app
COPY . .
RUN pip install django
EXPOSE 8000
CMD ["python", "manage.py", "runserver", "0.0.0.0:8000"]
EOF

# infra docker-compose
cat > infra/docker-compose.yaml <<'EOF'
version: "3.8"
services:
  redis:
    image: redis:7
    ports: ["6379:6379"]
  postgres:
    image: postgres:15
    environment:
      POSTGRES_PASSWORD: diam
    ports: ["5432:5432"]
  diameter:
    build: ../go-diameter-server
    ports: ["3868:3868", "8080:8080"]
    depends_on: [redis, postgres]
  control:
    build: ../django-control
    ports: ["8000:8000"]
    depends_on: [diameter]
EOF

# mkcerts
cat > infra/mkcerts.sh <<'EOF'
#!/usr/bin/env bash
set -e
# Very small helper to create local CA + a server cert for demos.
# For production use proper CA / cert-manager.
openssl genrsa -out ca.key 4096
openssl req -x509 -new -nodes -key ca.key -sha256 -days 3650 -subj "/CN=diameter-ca" -out ca.crt
openssl genrsa -out server.key 2048
openssl req -new -key server.key -subj "/CN=diameter-local.example" -out server.csr
openssl x509 -req -in server.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out server.crt -days 365 -sha256
echo "Generated ca.crt server.crt server.key in $(pwd)"
EOF
chmod +x infra/mkcerts.sh

# simple test client stub
cat > tests/integration/test_client.go <<'EOF'
package main
// Minimal test client stub (use go-diameter or net package to craft real messages)
import "fmt"
func main(){ fmt.Println("test client placeholder") }
EOF

echo "Scaffold created. Run 'cd diameter-platform/infra && docker-compose up --build' to start the demo stack."
