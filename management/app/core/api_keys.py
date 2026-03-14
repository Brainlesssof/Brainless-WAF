import secrets
import hashlib

def generate_api_key():
    """Generates a new API key with 'bwaf_live_' prefix."""
    random_part = secrets.token_urlsafe(32)
    return f"bwaf_live_{random_part}"

def hash_api_key(api_key: str):
    """Hashes the API key for safe storage."""
    return hashlib.sha256(api_key.encode()).hexdigest()
