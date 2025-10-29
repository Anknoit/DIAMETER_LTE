from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware
from fastapi.staticfiles import StaticFiles
from prometheus_client import start_http_server
from api import health, peers, simulate, sessions, messages, auth
from app.ui import routes as ui_routes
import os

app = FastAPI(title="Diameter Control API", version="1.0.0")

app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_methods=["GET","POST","PUT","DELETE","OPTIONS"],
    allow_headers=["*"],
)

# Serve static files (CSS, JS)
app.mount("/static", StaticFiles(directory="static"), name="static")

# Include API routes
app.include_router(auth.router, prefix="/api/v1")
app.include_router(health.router, prefix="/api/v1")
app.include_router(peers.router, prefix="/api/v1")
app.include_router(sessions.router, prefix="/api/v1")
app.include_router(simulate.router, prefix="/api/v1")
app.include_router(messages.router, prefix="/api/v1")

# Include UI routes
app.include_router(ui_routes.router)

# Start Prometheus metrics server
PROM_PORT = int(os.getenv("PROM_PORT", "9180"))
start_http_server(PROM_PORT)

@app.get("/api")
async def root():
    return {"service": "diameter-control", "version": "1.0"}
