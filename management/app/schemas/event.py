from pydantic import BaseModel
from typing import Optional, List, Any
from datetime import datetime

class EventBase(BaseModel):
    request_id: str
    timestamp: datetime
    src_ip: str
    method: str
    path: str
    action: str
    severity: str
    anomaly_score: int
    triggered_rules: List[Any]
    user_agent: str

class Event(EventBase):
    id: int

    class Config:
        from_attributes = True

class EventList(BaseModel):
    total: int
    items: List[Event]

class Stats(BaseModel):
    period: str
    requests_total: int
    requests_blocked: int
    block_rate_pct: float
    latency_p50_ms: float
    latency_p95_ms: float
    latency_p99_ms: float
    top_attack_types: List[dict]
    top_blocked_ips: List[dict]
