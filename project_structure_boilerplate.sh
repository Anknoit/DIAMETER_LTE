#!/usr/bin/env bash
set -e
ROOT="diameter-platform-fastapi"
mkdir -p "$ROOT"
cd "$ROOT"

# top-level dirs
mkdir -p go-diameter-server/cmd/server
mkdir -p go-diameter-server/internal/diameter
mkdir -p go-diameter-server/internal/api
mkdir -p go-diameter-server/internal/store
mkdir -p go-diameter-server/internal/config
mkdir -p go-diameter-server/internal/metrics
mkdir -p go-diameter-server/config

mkdir -p fastapi-control/app/api
mkdir -p fastapi-control/app
mkdir -p fastapi-control/scripts

mkdir -p infra
mkdir -p tests/integration

# .gitignore
cat > .gitignore <<'EOF'
# Go
/bin/
/vendor/
*.exe
*.log

# Python
__pycache__/
*.pyc
.env

# Docker
docker-compose.override.yml

# IDEs
.vscode/
.idea/
EOF

# README
cat > README.md <<'EOF'
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
EOF

######## GO SERVER SCAFFOLD ########

cat > go-diameter-server/cmd/server/main.go <<'EOF'
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"diam/internal/api"
)

func main() {
	log.Println("diameter server scaffold starting")

	// TODO: replace with real config loader
	go func() {
		mux := http.NewServeMux()
		api.RegisterHandlers(mux)
		log.Println("starting management API on :8080")
		if err := http.ListenAndServe(":8080", mux); err != nil {
			log.Fatalf("management server error: %v", err)
		}
	}()

	// Placeholder: real Diameter listener (TCP/TLS) to be implemented in internal/diameter
	log.Println("Diameter protocol engine not implemented in scaffold (see internal/diameter)")

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	<-ctx.Done()
	log.Println("shutting down diameter server scaffold")
	// graceful shutdown logic would go here
}
EOF

cat > go-diameter-server/internal/api/handlers.go <<'EOF'
package api

import (
	"encoding/json"
	"net/http"
)

// simple management API used by FastAPI control plane for demo
func RegisterHandlers(mux *http.ServeMux) {
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte(`{"status":"ok"}`))
	})

	mux.HandleFunc("/peers", func(w http.ResponseWriter, r *http.Request) {
		// TODO: implement persistent peers store and runtime add/remove
		if r.Method == http.MethodGet {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode([]map[string]interface{}{})
			return
		}
		w.WriteHeader(http.StatusNotImplemented)
	})

	mux.HandleFunc("/simulate", func(w http.ResponseWriter, r *http.Request) {
		// For demo: accept posted JSON and echo a fake answer
		var body map[string]interface{}
		_ = json.NewDecoder(r.Body).Decode(&body)
		resp := map[string]interface{}{
			"result": "ok (simulated)", "request": body,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	})
}
EOF

cat > go-diameter-server/internal/diameter/handler.go <<'EOF'
package diameter

// Protocol handlers placeholder
// Implement CER/CEA, DWR/DWA, ACR/ACA, AAR/AAA, CCR/CCA here using go-diameter.
func InitHandlers() {
	// register handlers with the chosen Diameter library
}
EOF

cat > go-diameter-server/Dockerfile <<'EOF'
FROM golang:1.20-alpine AS build
WORKDIR /src
COPY . .
RUN apk add --no-cache git && go env -w GOPATH=/go && go build -o /bin/diameter ./cmd/server

FROM alpine:3.18
COPY --from=build /bin/diameter /usr/local/bin/diameter
COPY config /etc/diameter
EXPOSE 3868 8080
CMD ["/usr/local/bin/diameter"]
EOF

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

######## FASTAPI CONTROL SCAFFOLD ########

cat > fastapi-control/requirements.txt <<'EOF'
fastapi
uvicorn[standard]
httpx
pydantic
python-jose
passlib[bcrypt]
aioredis
asyncpg
prometheus-client
python-multipart
EOF

cat > fastapi-control/app/main.py <<'EOF'
from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware
from app.api import health, peers, simulate, sessions, messages, auth
from prometheus_client import start_http_server
import os

app = FastAPI(title="Diameter Control API", version="1.0.0")

app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_methods=["GET","POST","PUT","DELETE","OPTIONS"],
    allow_headers=["*"],
)

app.include_router(auth.router, prefix="/api/v1")
app.include_router(health.router, prefix="/api/v1")
app.include_router(peers.router, prefix="/api/v1")
app.include_router(sessions.router, prefix="/api/v1")
app.include_router(simulate.router, prefix="/api/v1")
app.include_router(messages.router, prefix="/api/v1")

PROM_PORT = int(os.getenv("PROM_PORT", "9180"))
start_http_server(PROM_PORT)

@app.get("/")
async def root():
    return {"service": "diameter-control", "version": "1.0"}
EOF

cat > fastapi-control/app/api/health.py <<'EOF'
from fastapi import APIRouter
router = APIRouter()

@router.get("/health")
async def health():
    return {"status": "ok"}
EOF

cat > fastapi-control/app/api/auth.py <<'EOF'
from fastapi import APIRouter, Depends, HTTPException, status
from fastapi.security import OAuth2PasswordBearer, OAuth2PasswordRequestForm
from pydantic import BaseModel
from jose import jwt
from datetime import datetime, timedelta
from passlib.hash import bcrypt

router = APIRouter()
oauth2_scheme = OAuth2PasswordBearer(tokenUrl="/api/v1/token")

# demo user store - replace with DB
USERS = {"admin": {"username": "admin", "password_hash": bcrypt.hash("changeme"), "role": "admin"}}
SECRET_KEY = "replace-with-secure-secret"
ALGORITHM = "HS256"
ACCESS_TOKEN_EXPIRE_MINUTES = 60

class Token(BaseModel):
    access_token: str
    token_type: str

@router.post("/token", response_model=Token)
async def login(form_data: OAuth2PasswordRequestForm = Depends()):
    user = USERS.get(form_data.username)
    if not user or not bcrypt.verify(form_data.password, user["password_hash"]):
        raise HTTPException(status_code=status.HTTP_401_UNAUTHORIZED, detail="Invalid credentials")
    expire = datetime.utcnow() + timedelta(minutes=ACCESS_TOKEN_EXPIRE_MINUTES)
    payload = {"sub": user["username"], "exp": expire}
    token = jwt.encode(payload, SECRET_KEY, algorithm=ALGORITHM)
    return {"access_token": token, "token_type": "bearer"}

async def get_current_user(token: str = Depends(oauth2_scheme)):
    try:
        data = jwt.decode(token, SECRET_KEY, algorithms=[ALGORITHM])
        username = data.get("sub")
        if username is None:
            raise Exception()
        return USERS.get(username)
    except Exception:
        raise HTTPException(status_code=401, detail="Invalid token")
EOF

cat > fastapi-control/app/api/peers.py <<'EOF'
from fastapi import APIRouter, Depends, HTTPException
from pydantic import BaseModel
from typing import List
from .auth import get_current_user

router = APIRouter()
PEERS = {}

class PeerIn(BaseModel):
    id: str
    host: str
    ip: str
    port: int = 3868
    realm: str | None = None
    tls: bool = True
    ca_cert: str | None = None

class PeerOut(PeerIn):
    status: str = "unknown"
    last_seen: str | None = None

@router.get("/peers", response_model=List[PeerOut])
async def list_peers(user=Depends(get_current_user)):
    return [PeerOut(**p) for p in PEERS.values()]

@router.post("/peers", response_model=PeerOut)
async def create_peer(peer: PeerIn, user=Depends(get_current_user)):
    if peer.id in PEERS:
        raise HTTPException(status_code=409, detail="peer exists")
    PEERS[peer.id] = peer.dict()
    # TODO: call Go server management API to add peer, validate cert
    return PeerOut(**PEERS[peer.id])
EOF

cat > fastapi-control/app/api/simulate.py <<'EOF'
from fastapi import APIRouter, Depends
from pydantic import BaseModel
from .auth import get_current_user
from typing import Dict, Any
import httpx, os

router = APIRouter()

class SimRequest(BaseModel):
    type: str
    subtype: str | None = None
    session_id: str | None = None
    user: str | None = None
    peer_id: str | None = None
    avps: Dict[str, Any] = {}

GO_API = os.getenv("GO_API", "http://diameter:8080")

@router.post("/simulate")
async def simulate(req: SimRequest, user=Depends(get_current_user)):
    async with httpx.AsyncClient() as cl:
        r = await cl.post(f"{GO_API}/simulate", json=req.dict(), timeout=15.0)
        r.raise_for_status()
        return r.json()
EOF

cat > fastapi-control/app/api/sessions.py <<'EOF'
from fastapi import APIRouter, Depends
from .auth import get_current_user
router = APIRouter()

@router.get("/sessions")
async def list_sessions(user=Depends(get_current_user)):
    # TODO: call Go server or Redis to return active sessions
    return {"sessions": []}
EOF

cat > fastapi-control/app/api/messages.py <<'EOF'
from fastapi import APIRouter, WebSocket, WebSocketDisconnect
from typing import List
router = APIRouter()
clients: List[WebSocket] = []

@router.websocket("/ws/messages")
async def messages_ws(ws: WebSocket):
    await ws.accept()
    clients.append(ws)
    try:
        while True:
            # keep connection alive; real server pushes messages from Redis/Kafka
            await ws.receive_text()
    except WebSocketDisconnect:
        clients.remove(ws)
EOF

cat > fastapi-control/Dockerfile <<'EOF'
FROM python:3.11-slim
WORKDIR /app
COPY . .
RUN pip install --no-cache-dir --upgrade pip \
    && pip install -r requirements.txt
ENV PYTHONPATH=/app
EXPOSE 8000
CMD ["uvicorn", "app.main:app", "--host", "0.0.0.0", "--port", "8000"]
EOF

cat > fastapi-control/.env <<'EOF'
GO_API=http://diameter:8080
PROM_PORT=9180
EOF

######## INFRA (docker-compose + mkcerts) ########

cat > infra/docker-compose.yaml <<'EOF'
version: "3.8"
services:
  redis:
    image: redis:7
    ports:
      - "6379:6379"

  postgres:
    image: postgres:15
    environment:
      POSTGRES_USER: diam
      POSTGRES_PASSWORD: diam
      POSTGRES_DB: diameter
    ports:
      - "5432:5432"

  diameter:
    build: ../go-diameter-server
    ports:
      - "3868:3868"
      - "8080:8080"
    depends_on:
      - redis
      - postgres

  control:
    build: ../fastapi-control
    environment:
      - GO_API=http://diameter:8080
      - PROM_PORT=9180
    ports:
      - "8000:8000"
    depends_on:
      - diameter
EOF

cat > infra/mkcerts.sh <<'EOF'
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
EOF
chmod +x infra/mkcerts.sh

######## TEST CLIENT ########

cat > tests/integration/test_client.go <<'EOF'
package main

import "fmt"
func main() {
    fmt.Println("Integration test client placeholder. Use go-diameter library or httpx against management API.")
}
EOF

echo "Scaffold created in $(pwd). To start the stack:"
echo "  cd infra"
echo "  ./mkcerts.sh"
echo "  docker-compose up --build"
