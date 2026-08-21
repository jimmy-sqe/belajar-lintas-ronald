from __future__ import annotations

import contextlib
import os
from collections.abc import Generator
from pathlib import Path
from typing import Any

import pytest

from framework.config import PROJECTS, Settings

pytest_plugins = [
    "tests.fixtures.browser",
    "tests.fixtures.auth",
    # boilerplate:axis=api-testing option=wired START
    "tests.fixtures.api",
    "tests.fixtures.factories",
    # boilerplate:axis=api-testing option=wired END
]


def pytest_configure(config: pytest.Config) -> None:
    # pytest-html writes report.html at sessionfinish; pytest-playwright may
    # not create the --output dir until a test produces an artifact. allure-pytest
    # writes per-test JSON files synchronously and errors if --alluredir is missing.
    # Pre-create all three so none of the reporter plugins races.
    output = config.getoption("--output", default="pytest-results")
    Path(output).mkdir(parents=True, exist_ok=True)
    html_path = config.getoption("--html", default=None)
    if html_path:
        Path(html_path).parent.mkdir(parents=True, exist_ok=True)
    allure_dir = config.getoption("--alluredir", default=None)
    if allure_dir:
        Path(allure_dir).mkdir(parents=True, exist_ok=True)


@pytest.hookimpl(trylast=True)
def pytest_sessionfinish(session: pytest.Session) -> None:
    # Idempotent re-mkdir: pytest-playwright cleans empty artifact subdirs at
    # session end, which can rmdir() the parent if no failures occurred. Pre-html-write.
    output = session.config.getoption("--output", default="pytest-results")
    Path(output).mkdir(parents=True, exist_ok=True)
    html_path = session.config.getoption("--html", default=None)
    if html_path:
        Path(html_path).parent.mkdir(parents=True, exist_ok=True)
    allure_dir = session.config.getoption("--alluredir", default=None)
    if allure_dir:
        Path(allure_dir).mkdir(parents=True, exist_ok=True)


def _active_project(request: pytest.FixtureRequest) -> str | None:
    for marker in PROJECTS:
        if request.node.get_closest_marker(marker):
            return marker
    return None


@pytest.fixture(scope="session")
def settings() -> Settings:
    return Settings()


@pytest.fixture
def base_url(request: pytest.FixtureRequest, settings: Settings) -> str:
    project = _active_project(request)
    if project:
        return settings.project(project).base_url
    return settings.base_url


@pytest.fixture
def credentials(request: pytest.FixtureRequest, settings: Settings) -> dict[str, str | None]:
    project = _active_project(request)
    if project:
        cfg = settings.project(project)
        return {"username": cfg.username, "password": cfg.password}
    return {"username": settings.test_username, "password": settings.test_password}


def pytest_collection_modifyitems(config: pytest.Config, items: list[pytest.Item]) -> None:
    for item in items:
        for part in Path(str(item.fspath)).parts:
            if part in PROJECTS:
                item.add_marker(getattr(pytest.mark, part))
                break
        if item.get_closest_marker("flaky"):
            item.add_marker(pytest.mark.flaky(reruns=2, reruns_delay=2))


def _artifact_dir_for(item: pytest.Item) -> Path:
    output = Path(item.config.getoption("--output", default="pytest-results"))
    worker = os.environ.get("PYTEST_XDIST_WORKER", "gw0")
    safe_nodeid = (
        item.nodeid.replace("::", "__").replace("/", "_").replace("[", "_").replace("]", "")
    )
    project_part = "unknown"
    for part in Path(str(item.fspath)).parts:
        if part in PROJECTS:
            project_part = part
            break
    d = output / project_part / f"{safe_nodeid}__{worker}"
    d.mkdir(parents=True, exist_ok=True)
    return d


@pytest.hookimpl(hookwrapper=True)
def pytest_runtest_makereport(
    item: pytest.Item, call: pytest.CallInfo[None]
) -> Generator[None, Any, None]:
    outcome = yield
    rep = outcome.get_result()
    if rep.when != "call" or not rep.failed:
        return
    funcargs: dict[str, Any] = getattr(item, "funcargs", {}) or {}
    page = funcargs.get("logged_in_page") or funcargs.get("page")
    if page is None:
        return
    artifact_dir = _artifact_dir_for(item)
    png = artifact_dir / "failure.png"
    try:
        page.screenshot(path=str(png), full_page=True)
    except Exception:
        return
    url = ""
    with contextlib.suppress(Exception):
        url = page.url

    extras = getattr(rep, "extras", []) or []
    try:
        from pytest_html import extras as html_extras

        extras.append(html_extras.image(str(png)))
        if url:
            extras.append(html_extras.text(url, name="URL at failure"))
    except Exception:
        pass
    rep.extras = extras

    # boilerplate:axis=reporting option=allure START
    # Allure attachments — mirror what pytest-html gets so both reports stay in sync.
    try:
        import allure

        allure.attach.file(
            str(png),
            name="failure.png",
            attachment_type=allure.attachment_type.PNG,
        )
        if url:
            allure.attach(url, name="URL at failure", attachment_type=allure.attachment_type.TEXT)
        # If pytest-playwright produced a trace, attach the path so Allure links to it.
        trace_zip = artifact_dir / "trace.zip"
        if trace_zip.exists():
            allure.attach.file(
                str(trace_zip),
                name="playwright-trace.zip",
                extension="zip",
            )
    except Exception:
        pass
    # boilerplate:axis=reporting option=allure END
