from datetime import datetime
from collections import defaultdict
from typing import Dict, Any

# Track hit counts and request metadata
request_stats: Dict[str, Any] = {
    "endpoints": defaultdict(lambda: {"count": 0, "requests": []})
}


def log_request(path: str, method: str, params: Dict[str, Any]):
    """Record a request hit and relevant info."""
    entry = request_stats["endpoints"][path]
    entry["count"] += 1
    entry["requests"].append({
        "method": method,
        "params": params,
        "timestamp": datetime.utcnow().isoformat() + "Z"
    })


def get_request_stats() -> Dict[str, Any]:
    """Return a snapshot of current stats."""
    return request_stats
