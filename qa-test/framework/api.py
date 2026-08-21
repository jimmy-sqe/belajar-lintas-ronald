from __future__ import annotations

import requests

from framework.api_clients.base import BaseAPIClient
from framework.config import Settings

# Populated by /generate-api-tests, which writes the per-project client into
# framework/api_clients/ and registers it here.
API_CLIENTS: dict[str, type[BaseAPIClient]] = {}


def api_for(project: str, settings: Settings) -> BaseAPIClient:
    """Build a logged-in API client for the given project marker."""
    if project not in API_CLIENTS:
        raise KeyError(f"No APIClient registered for project {project!r}")
    cfg = settings.project(project)
    cls = API_CLIENTS[project]
    client = cls(cfg.base_url, requests.Session())
    client.login(cfg.username, cfg.password)
    return client
