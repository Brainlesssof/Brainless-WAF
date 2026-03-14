from typing import Any, Dict
from fastapi import APIRouter, Depends, HTTPException, status
from sqlalchemy.orm import Session
from app.api.v1 import deps
from app.db.session import get_db
from app.models import models
from app.schemas import config as config_schema
from app.services.notifier import get_notifier, WafNotifier

router = APIRouter()

@router.get("/config", response_model=Dict[str, Any])
async def read_config(
    db: Session = Depends(get_db),
    current_user: models.User = Depends(deps.get_current_active_user)
) -> Any:
    configs = db.query(models.Config).all()
    return {c.key: c.value for c in configs}

@router.get("/config/{key}", response_model=config_schema.Config)
async def read_config_item(
    key: str,
    db: Session = Depends(get_db),
    current_user: models.User = Depends(deps.get_current_active_user)
) -> Any:
    config = db.query(models.Config).filter(models.Config.key == key).first()
    if not config:
        raise HTTPException(status_code=404, detail="Config key not found")
    return config

@router.patch("/config", response_model=Dict[str, Any])
async def update_config(
    config_in: Dict[str, Any],
    db: Session = Depends(get_db),
    current_user: models.User = Depends(deps.check_admin),
    notifier: WafNotifier = Depends(get_notifier)
) -> Any:
    for key, value in config_in.items():
        db_config = db.query(models.Config).filter(models.Config.key == key).first()
        if db_config:
            db_config.value = value
        else:
            db_config = models.Config(key=key, value=value)
            db.add(db_config)
        await notifier.notify_config_change(key, value)
    
    db.commit()
    # Return updated config
    configs = db.query(models.Config).all()
    return {c.key: c.value for c in configs}
