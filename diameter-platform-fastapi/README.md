# Diameter Platform (FastAPI control plane scaffold)

Scaffold includes:
- go-diameter-server (Go scaffold for protocol + management API)
- fastapi-control (FastAPI control plane for operator UI & orchestration)
- infra/docker-compose.yaml (local dev stack)
- infra/mkcerts.sh (demo cert creation)

Run locally:
  cd infra
  docker-compose up --build

FastAPI docs: http://localhost:8000/docs
Go server health: http://localhost:8080/health
