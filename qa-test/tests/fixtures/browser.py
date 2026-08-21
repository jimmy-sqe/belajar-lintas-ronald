from __future__ import annotations

from typing import Any

import pytest

from framework.config import Settings


def _project_marker(request: pytest.FixtureRequest) -> str | None:
    from framework.config import PROJECTS

    for marker in PROJECTS:
        if request.node.get_closest_marker(marker):
            return marker
    return None


def _resolve_viewport(request: pytest.FixtureRequest, settings: Settings) -> dict[str, int]:
    marker = request.node.get_closest_marker("viewport")
    if marker and len(marker.args) >= 2:
        return {"width": int(marker.args[0]), "height": int(marker.args[1])}
    project = _project_marker(request)
    if project:
        cfg = settings.project(project)
        return {"width": cfg.viewport[0], "height": cfg.viewport[1]}
    return {"width": 1440, "height": 900}


@pytest.fixture
def browser_context_args(
    request: pytest.FixtureRequest,
    settings: Settings,
    pytestconfig: pytest.Config,
) -> dict[str, Any]:
    args: dict[str, Any] = {
        "viewport": _resolve_viewport(request, settings),
        "locale": "en-US",
        "ignore_https_errors": True,
    }
    project = _project_marker(request)
    if project:
        cfg = settings.project(project)
        if cfg.base_url:
            args["base_url"] = cfg.base_url
    elif settings.base_url:
        args["base_url"] = settings.base_url

    video = pytestconfig.getoption("--video", default="off")
    if video in ("on", "retain-on-failure"):
        args["record_video_dir"] = pytestconfig.getoption("--output", default="test-results")
    return args
