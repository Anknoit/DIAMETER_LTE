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
