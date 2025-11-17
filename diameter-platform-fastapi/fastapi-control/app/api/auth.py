# app/api/auth.py
from fastapi import APIRouter, Depends, HTTPException, status
from fastapi.security import OAuth2PasswordBearer, OAuth2PasswordRequestForm
from pydantic import BaseModel
from jose import jwt
from datetime import datetime, timedelta
from passlib.context import CryptContext
from typing import Dict, Optional

router = APIRouter()
oauth2_scheme = OAuth2PasswordBearer(tokenUrl="/api/v1/token")

# CryptContext - recommended way to handle hashing
pwd_context = CryptContext(schemes=["bcrypt"], deprecated="auto")

SECRET_KEY = "replace-with-secure-secret"
ALGORITHM = "HS256"
ACCESS_TOKEN_EXPIRE_MINUTES = 60

class Token(BaseModel):
    access_token: str
    token_type: str

# persistent user store should be DB; this is an in-memory demo store.
# For demo we store a precomputed password hash or generate it lazily.
USERS: Dict[str, Dict[str, Optional[str]]] = {
    # Example entry format:
    # "admin": {"username": "admin", "password_hash": "$2b$12$....", "role": "admin"}
}

def create_demo_user_once(username: str = "admin", raw_password: str = "changeme"):
    """
    If the demo user is not present, create it and store the hash.
    Avoid hashing at module import time for safety in certain environments.
    """
    if username in USERS and USERS[username].get("password_hash"):
        return
    # bcrypt will raise if password > 72 bytes; we proactively truncate
    safe_pw = raw_password if len(raw_password.encode("utf-8")) <= 72 else raw_password.encode("utf-8")[:72].decode("utf-8", errors="ignore")
    USERS[username] = {
        "username": username,
        "password_hash": pwd_context.hash(safe_pw),
        "role": "admin"
    }

# Create demo user lazily the first time an auth route is called.
# Optionally you may call create_demo_user_once() from startup event.
create_demo_user_once()

@router.post("/token", response_model=Token)
async def login(form_data: OAuth2PasswordRequestForm = Depends()):
    user = USERS.get(form_data.username)
    if not user:
        raise HTTPException(status_code=status.HTTP_401_UNAUTHORIZED, detail="Invalid credentials")
    # ensure we truncate long passwords in verification path too
    raw_pw = form_data.password
    if len(raw_pw.encode("utf-8")) > 72:
        raw_pw = raw_pw.encode("utf-8")[:72].decode("utf-8", errors="ignore")

    if not pwd_context.verify(raw_pw, user["password_hash"]):
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
        user = USERS.get(username)
        if not user:
            raise Exception()
        return user
    except Exception:
        raise HTTPException(status_code=401, detail="Invalid token")
