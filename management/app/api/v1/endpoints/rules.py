from typing import Any, List, Optional
from fastapi import APIRouter, Depends, HTTPException, status, Query
from sqlalchemy.orm import Session
from app.api.v1 import deps
from app.db.session import get_db
from app.models import models
from app.schemas import rule as rule_schema
from app.services.notifier import get_notifier, WafNotifier

router = APIRouter()

@router.get("/rules", response_model=rule_schema.RuleList)
async def read_rules(
    db: Session = Depends(get_db),
    skip: int = 0,
    limit: int = 100,
    enabled: Optional[bool] = None,
    current_user: models.User = Depends(deps.get_current_active_user)
) -> Any:
    query = db.query(models.Rule)
    if enabled is not None:
        query = query.filter(models.Rule.enabled == enabled)
    
    total = query.count()
    rules = query.offset(skip).limit(limit).all()
    
    return {"total": total, "items": rules}

@router.post("/rules", response_model=rule_schema.Rule, status_code=status.HTTP_201_CREATED)
async def create_rule(
    *,
    db: Session = Depends(get_db),
    rule_in: rule_schema.RuleCreate,
    current_user: models.User = Depends(deps.check_admin),
    notifier: WafNotifier = Depends(get_notifier)
) -> Any:
    db_rule = models.Rule(**rule_in.dict())
    db.add(db_rule)
    try:
        db.commit()
        db.refresh(db_rule)
        await notifier.notify_rule_change([db_rule.rule_id])
    except Exception as e:
        db.rollback()
        raise HTTPException(status_code=400, detail="Rule ID already exists or invalid data")
    return db_rule

@router.get("/rules/{id}", response_model=rule_schema.Rule)
async def read_rule(
    *,
    db: Session = Depends(get_db),
    id: int,
    current_user: models.User = Depends(deps.get_current_active_user)
) -> Any:
    rule = db.query(models.Rule).filter(models.Rule.id == id).first()
    if not rule:
        raise HTTPException(status_code=404, detail="Rule not found")
    return rule

@router.put("/rules/{id}", response_model=rule_schema.Rule)
async def update_rule(
    *,
    db: Session = Depends(get_db),
    id: int,
    rule_in: rule_schema.RuleUpdate,
    current_user: models.User = Depends(deps.check_admin),
    notifier: WafNotifier = Depends(get_notifier)
) -> Any:
    rule = db.query(models.Rule).filter(models.Rule.id == id).first()
    if not rule:
        raise HTTPException(status_code=404, detail="Rule not found")
    
    update_data = rule_in.dict(exclude_unset=True)
    for field, value in update_data.items():
        setattr(rule, field, value)
    
    db.add(rule)
    db.commit()
    db.refresh(rule)
    await notifier.notify_rule_change([rule.rule_id])
    return rule

@router.delete("/rules/{id}", status_code=status.HTTP_204_NO_CONTENT)
async def delete_rule(
    *,
    db: Session = Depends(get_db),
    id: int,
    current_user: models.User = Depends(deps.check_admin),
    notifier: WafNotifier = Depends(get_notifier)
) -> Any:
    rule = db.query(models.Rule).filter(models.Rule.id == id).first()
    if not rule:
        raise HTTPException(status_code=404, detail="Rule not found")
    if rule.source == "crs":
        raise HTTPException(status_code=400, detail="Built-in rules cannot be deleted, only disabled")
    
    rule_id = rule.rule_id
    db.delete(rule)
    db.commit()
    await notifier.notify_rule_change([rule_id])
    return None
