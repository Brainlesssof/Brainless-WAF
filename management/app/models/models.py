from sqlalchemy import Column, Integer, String, Boolean, DateTime, JSON, ForeignKey, Text
from sqlalchemy.sql import func
from app.db.session import Base

class User(Base):
    __tablename__ = "users"

    id = Column(Integer, primary_key=True, index=True)
    username = Column(String, unique=True, index=True, nullable=False)
    email = Column(String, unique=True, index=True, nullable=False)
    hashed_password = Column(String, nullable=False)
    is_active = Column(Boolean, default=True)
    role = Column(String, default="analyst")  # admin, analyst, readonly
    created_at = Column(DateTime(timezone=True), server_default=func.now())

class Rule(Base):
    __tablename__ = "rules"

    id = Column(Integer, primary_key=True, index=True)
    rule_id = Column(Integer, unique=True, index=True, nullable=False)
    rule_text = Column(Text, nullable=False)
    description = Column(String)
    enabled = Column(Boolean, default=True)
    source = Column(String, default="custom")  # crs, custom
    phase = Column(Integer, default=2)
    severity = Column(String, default="NOTICE")
    created_at = Column(DateTime(timezone=True), server_default=func.now())
    updated_at = Column(DateTime(timezone=True), onupdate=func.now())

class IPList(Base):
    __tablename__ = "ip_lists"

    id = Column(Integer, primary_key=True, index=True)
    cidr = Column(String, index=True, nullable=False)
    list_type = Column(String, nullable=False)  # blocklist, allowlist
    reason = Column(String)
    expires_at = Column(DateTime(timezone=True), nullable=True)
    created_at = Column(DateTime(timezone=True), server_default=func.now())

class Config(Base):
    __tablename__ = "config"

    id = Column(Integer, primary_key=True, index=True)
    key = Column(String, unique=True, index=True, nullable=False)
    value = Column(JSON, nullable=False)
    updated_at = Column(DateTime(timezone=True), onupdate=func.now())

class Event(Base):
    __tablename__ = "events"

    id = Column(Integer, primary_key=True, index=True)
    request_id = Column(String, index=True)
    timestamp = Column(DateTime(timezone=True), server_default=func.now())
    src_ip = Column(String, index=True)
    method = Column(String)
    path = Column(String)
    action = Column(String)  # block, detect, allow
    severity = Column(String, index=True)
    anomaly_score = Column(Integer)
    triggered_rules = Column(JSON)  # List of rule objects
    user_agent = Column(String)
