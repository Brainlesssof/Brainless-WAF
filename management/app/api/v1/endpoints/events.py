from typing import Any, List, Optional
from fastapi import APIRouter, Depends, HTTPException, Query
from sqlalchemy.orm import Session
from datetime import datetime
from app.api.v1 import deps
from app.db.session import get_db
from app.models import models
from app.schemas import event as event_schema

router = APIRouter()

@router.get("/events", response_model=event_schema.EventList)
async def read_events(
    db: Session = Depends(get_db),
    skip: int = 0,
    limit: int = 50,
    severity: Optional[str] = None,
    action: Optional[str] = None,
    src_ip: Optional[str] = None,
    current_user: models.User = Depends(deps.get_current_active_user)
) -> Any:
    query = db.query(models.Event)
    if severity:
        query = query.filter(models.Event.severity == severity)
    if action:
        query = query.filter(models.Event.action == action)
    if src_ip:
        query = query.filter(models.Event.src_ip == src_ip)
    
    total = query.count()
    items = query.order_by(models.Event.timestamp.desc()).offset(skip).limit(limit).all()
    return {"total": total, "items": items}

@router.get("/stats", response_model=event_schema.Stats)
async def read_stats(
    period: str = "24h",
    db: Session = Depends(get_db),
    current_user: models.User = Depends(deps.get_current_active_user)
) -> Any:
    # This would typically involve complex aggregation/caching
    # Returning mock data for now as per api.md examples
    return {
        "period": period,
        "requests_total": 1482930,
        "requests_blocked": 342,
        "block_rate_pct": 0.023,
        "latency_p50_ms": 0.8,
        "latency_p95_ms": 2.1,
        "latency_p99_ms": 4.7,
        "top_attack_types": [
            {"type": "SQLi", "count": 156},
            {"type": "XSS", "count": 89},
            {"type": "PathTraversal", "count": 43}
        ],
        "top_blocked_ips": [
            {"ip": "203.0.113.42", "count": 87}
        ]
    }
