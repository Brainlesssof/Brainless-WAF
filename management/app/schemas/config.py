from pydantic import BaseModel
from typing import Any, Dict

class ConfigBase(BaseModel):
    key: str
    value: Any

class ConfigUpdate(BaseModel):
    value: Any

class Config(ConfigBase):
    class Config:
        from_attributes = True

class ConfigPatch(BaseModel):
    # For bulk updates if needed
    updates: Dict[str, Any]
