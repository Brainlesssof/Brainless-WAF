from pydantic import BaseModel
from typing import Optional, List
from datetime import datetime

class IPListBase(BaseModel):
    cidr: str
    list_type: str  # blocklist, allowlist
    reason: Optional[str] = None
    expires_at: Optional[datetime] = None

class IPListCreate(IPListBase):
    pass

class IPList(IPListBase):
    id: int
    created_at: datetime

    class Config:
        from_attributes = True

class IPListResponse(BaseModel):
    total: int
    items: List[IPList]
