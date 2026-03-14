from fastapi import APIRouter

router = APIRouter()

@router.get("/health")
async def health_check():
    return {
        "status": "ok",
        "version": "0.1.0",
        "uptime_seconds": 0,  # Placeholder
        "rules_loaded": 0,    # Placeholder
        "mode": "learning",
        "upstream_healthy": True
    }
