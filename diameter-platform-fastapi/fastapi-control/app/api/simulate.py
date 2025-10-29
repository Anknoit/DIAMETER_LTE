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
