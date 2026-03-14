import pytest
from fastapi.testclient import TestClient
from sqlalchemy import create_engine
from sqlalchemy.orm import sessionmaker
from app.main import app
from app.db.session import Base, get_db
from app.core import security
from app.models import models

# Use SQLite for testing
SQLALCHEMY_DATABASE_URL = "sqlite:///./test.db"
engine = create_engine(SQLALCHEMY_DATABASE_URL, connect_args={"check_same_thread": False})
TestingSessionLocal = sessionmaker(autocommit=False, autoflush=False, bind=engine)

def override_get_db():
    try:
        db = TestingSessionLocal()
        yield db
    finally:
        db.close()

app.dependency_overrides[get_db] = override_get_db

client = TestClient(app)

@pytest.fixture(autouse=True)
def setup_db():
    Base.metadata.create_all(bind=engine)
    db = TestingSessionLocal()
    # Create test admin user
    hashed_password = security.get_password_hash("testpass")
    user = models.User(username="testadmin", email="admin@example.com", hashed_password=hashed_password, role="admin")
    db.add(user)
    db.commit()
    yield
    Base.metadata.drop_all(bind=engine)

def test_health_check():
    response = client.get("/api/v1/health")
    assert response.status_code == 200
    assert response.json()["status"] == "ok"

def test_login_success():
    response = client.post(
        "/api/v1/auth/token",
        data={"username": "testadmin", "password": "testpass"}
    )
    assert response.status_code == 200
    assert "access_token" in response.json()

def test_login_failure():
    response = client.post(
        "/api/v1/auth/token",
        data={"username": "testadmin", "password": "wrongpassword"}
    )
    assert response.status_code == 401

def test_get_rules_unauthorized():
    response = client.get("/api/v1/rules")
    assert response.status_code == 401

def test_create_rule_success():
    # Login first
    login_response = client.post(
        "/api/v1/auth/token",
        data={"username": "testadmin", "password": "testpass"}
    )
    token = login_response.json()["access_token"]
    
    # Create rule
    response = client.post(
        "/api/v1/rules",
        json={
            "rule_id": 1001,
            "rule_text": "SecRule ARGS '@rx test' id:1001",
            "description": "Test rule",
            "severity": "CRITICAL"
        },
        headers={"Authorization": f"Bearer {token}"}
    )
    assert response.status_code == 201
    assert response.json()["rule_id"] == 1001
