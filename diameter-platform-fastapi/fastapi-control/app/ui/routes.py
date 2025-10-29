from fastapi import APIRouter, Request
from fastapi.templating import Jinja2Templates
from fastapi.responses import HTMLResponse

router = APIRouter()
templates = Jinja2Templates(directory="app/templates")

@router.get("/", response_class=HTMLResponse)
async def dashboard(request: Request):
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
