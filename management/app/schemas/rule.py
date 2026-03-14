from pydantic import BaseModel, Field
from typing import Optional, List
from datetime import datetime

class RuleBase(BaseModel):
    rule_id: int
    rule_text: str
    description: Optional[str] = None
    enabled: bool = True
    source: str = "custom"
    phase: int = 2
    severity: str = "NOTICE"

class RuleCreate(RuleBase):
    pass

class RuleUpdate(BaseModel):
    rule_text: Optional[str] = None
    description: Optional[str] = None
    enabled: Optional[bool] = None
    severity: Optional[str] = None

class Rule(RuleBase):
    id: int
    created_at: datetime
    updated_at: Optional[datetime] = None

    class Config:
        from_attributes = True

class RuleList(BaseModel):
    total: int
    items: List[Rule]
