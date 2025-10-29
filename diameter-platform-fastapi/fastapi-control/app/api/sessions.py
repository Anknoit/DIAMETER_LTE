from fastapi import APIRouter, Depends
from .auth import get_current_user
router = APIRouter()

@router.get("/sessions")
async def list_sessions(user=Depends(get_current_user)):
    # TODO: call Go server or Redis to return active sessions
    return {"sessions": []}
