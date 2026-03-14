from typing import Any, List, Optional
from fastapi import APIRouter, Depends, HTTPException, status
from sqlalchemy.orm import Session
from app.api.v1 import deps
from app.db.session import get_db
from app.models import models
from app.schemas import ip_list as ip_list_schema

router = APIRouter()

@router.get("/blocklist", response_model=ip_list_schema.IPListResponse)
async def read_blocklist(
    db: Session = Depends(get_db),
    skip: int = 0,
    limit: int = 100,
    current_user: models.User = Depends(deps.get_current_active_user)
) -> Any:
    query = db.query(models.IPList).filter(models.IPList.list_type == "blocklist")
    total = query.count()
    items = query.offset(skip).limit(limit).all()
    return {"total": total, "items": items}

@router.post("/blocklist", response_model=ip_list_schema.IPList, status_code=status.HTTP_201_CREATED)
async def add_to_blocklist(
    *,
    db: Session = Depends(get_db),
    item_in: ip_list_schema.IPListCreate,
    current_user: models.User = Depends(deps.check_admin)
) -> Any:
    item_in.list_type = "blocklist"
    db_item = models.IPList(**item_in.dict())
    db.add(db_item)
    db.commit()
    db.refresh(db_item)
    return db_item

@router.delete("/blocklist/{id}", status_code=status.HTTP_204_NO_CONTENT)
async def remove_from_blocklist(
    *,
    db: Session = Depends(get_db),
    id: int,
    current_user: models.User = Depends(deps.check_admin)
) -> Any:
    item = db.query(models.IPList).filter(models.IPList.id == id, models.IPList.list_type == "blocklist").first()
    if not item:
        raise HTTPException(status_code=404, detail="Entry not found in blocklist")
    db.delete(item)
    db.commit()
    return None

# Similar endpoints for allowlist
@router.get("/allowlist", response_model=ip_list_schema.IPListResponse)
async def read_allowlist(
    db: Session = Depends(get_db),
    skip: int = 0,
    limit: int = 100,
    current_user: models.User = Depends(deps.get_current_active_user)
) -> Any:
    query = db.query(models.IPList).filter(models.IPList.list_type == "allowlist")
    total = query.count()
    items = query.offset(skip).limit(limit).all()
    return {"total": total, "items": items}

@router.post("/allowlist", response_model=ip_list_schema.IPList, status_code=status.HTTP_201_CREATED)
async def add_to_allowlist(
    *,
    db: Session = Depends(get_db),
    item_in: ip_list_schema.IPListCreate,
    current_user: models.User = Depends(deps.check_admin)
) -> Any:
    item_in.list_type = "allowlist"
    db_item = models.IPList(**item_in.dict())
    db.add(db_item)
    db.commit()
    db.refresh(db_item)
    return db_item

@router.delete("/allowlist/{id}", status_code=status.HTTP_204_NO_CONTENT)
async def remove_from_allowlist(
    *,
    db: Session = Depends(get_db),
    id: int,
    current_user: models.User = Depends(deps.check_admin)
) -> Any:
    item = db.query(models.IPList).filter(models.IPList.id == id, models.IPList.list_type == "allowlist").first()
    if not item:
        raise HTTPException(status_code=404, detail="Entry not found in allowlist")
    db.delete(item)
    db.commit()
    return None
