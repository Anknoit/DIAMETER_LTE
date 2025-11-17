# app/ui/routes.py
from fastapi import APIRouter, Request
from fastapi.responses import HTMLResponse
from fastapi.templating import Jinja2Templates
from pathlib import Path
from datetime import datetime

router = APIRouter()

# compute templates directory relative to this file
BASE_DIR = Path(__file__).resolve().parent.parent  # points to <project>/app
TEMPLATES_DIR = BASE_DIR / "templates"
templates = Jinja2Templates(directory=str(TEMPLATES_DIR))

# expose helper functions / globals to all templates
templates.env.globals["now"] = datetime.utcnow
templates.env.globals["static_url"] = lambda path: f"/static/{path.lstrip('/')}"  # optional helper

@router.get("/", response_class=HTMLResponse)
async def dashboard(request: Request):
    # you can also pass per-request context here if needed
    return templates.TemplateResponse("dashboard.html", {"request": request, "title": "Dashboard"})

@router.get("/peers", response_class=HTMLResponse)
async def peers(request: Request):
    return templates.TemplateResponse("peers.html", {"request": request, "title": "Peers"})

@router.get("/sessions", response_class=HTMLResponse)
async def sessions(request: Request):
    return templates.TemplateResponse("sessions.html", {"request": request, "title": "Sessions"})

@router.get("/simulate", response_class=HTMLResponse)
async def simulate(request: Request):
    return templates.TemplateResponse("simulate.html", {"request": request, "title": "Simulator"})

@router.get("/messages", response_class=HTMLResponse)
async def messages(request: Request):
    return templates.TemplateResponse("messages.html", {"request": request, "title": "Messages"})

@router.get("/certs", response_class=HTMLResponse)
async def certs(request: Request):
    return templates.TemplateResponse("certs.html", {"request": request, "title": "Certificates"})
