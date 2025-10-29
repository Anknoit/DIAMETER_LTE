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
