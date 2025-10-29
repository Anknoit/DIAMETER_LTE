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
