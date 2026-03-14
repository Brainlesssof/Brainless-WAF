from fastapi import FastAPI
from app.api.v1.endpoints import health, auth, rules, config, ip_lists, events

app = FastAPI(
    title="Brainless WAF Management API",
    description="API for managing rules, configuration, and telemetry for Brainless WAF",
    version="0.1.0",
    docs_url="/api/v1/docs",
    redoc_url="/api/v1/redoc",
    openapi_url="/api/v1/openapi.json",
)

# Include routers
app.include_router(health.router, prefix="/api/v1", tags=["system"])
app.include_router(auth.router, prefix="/api/v1", tags=["auth"])
app.include_router(rules.router, prefix="/api/v1", tags=["rules"])
app.include_router(config.router, prefix="/api/v1", tags=["config"])
app.include_router(ip_lists.router, prefix="/api/v1", tags=["ip-lists"])
app.include_router(events.router, prefix="/api/v1", tags=["events"])

@app.get("/")
async def root():
    return {"message": "Welcome to Brainless WAF Management API. Visit /api/v1/docs for documentation."}
