from __future__ import annotations

from collections.abc import Iterator

import pytest

from framework.api import api_for
from framework.api_clients.base import BaseAPIClient
from framework.config import PROJECTS, Settings


def _active_project(request: pytest.FixtureRequest) -> str | None:
    for marker in PROJECTS:
        if request.node.get_closest_marker(marker):
            return marker
    return None


@pytest.fixture
def api(request: pytest.FixtureRequest, settings: Settings) -> Iterator[BaseAPIClient]:
    """Fresh API login per test. Requires a project marker on the test."""
    project = _active_project(request)
    if project is None:
        raise RuntimeError(
            "api fixture requires a project marker (e.g. @pytest.mark.belajar_lintas_ronald)"
        )
    client = api_for(project, settings)
    try:
        yield client
    finally:
        client.s.close()
