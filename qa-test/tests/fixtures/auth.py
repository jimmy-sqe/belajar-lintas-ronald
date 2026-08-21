from __future__ import annotations

import pytest
from playwright.sync_api import Page

from framework.auth import login_as
from framework.config import PROJECTS, Settings


def _active_project(request: pytest.FixtureRequest) -> str | None:
    for marker in PROJECTS:
        if request.node.get_closest_marker(marker):
            return marker
    return None


@pytest.fixture
def logged_in_page(page: Page, request: pytest.FixtureRequest, settings: Settings) -> Page:
    """Fresh UI login per test. Requires a project marker on the test."""
    project = _active_project(request)
    if project is None:
        raise RuntimeError(
            "logged_in_page requires a project marker (e.g. @pytest.mark.belajar_lintas_ronald)"
        )
    login_as(page, project, settings)
    return page
