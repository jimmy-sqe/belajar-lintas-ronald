from __future__ import annotations

from collections.abc import Iterator
from types import ModuleType

import pytest

from framework.config import PROJECTS
from framework.factories.base import CleanupRegistry

FACTORIES: dict[str, ModuleType] = {
    # Register a factory module per project as factories land, e.g.:
    # "belajar_lintas_ronald": belajar_lintas_ronald_factories,
}


def _active_project(request: pytest.FixtureRequest) -> str | None:
    for marker in PROJECTS:
        if request.node.get_closest_marker(marker):
            return marker
    return None


@pytest.fixture
def cleanup_registry() -> Iterator[CleanupRegistry]:
    """Per-test LIFO cleanup registry; drains on teardown.

    Factories use this to register a delete-after-test action when they
    create a resource. Use directly only when you need cleanup outside a
    factory call.
    """
    reg = CleanupRegistry()
    try:
        yield reg
    finally:
        reg.drain()


@pytest.fixture
def factory(request: pytest.FixtureRequest) -> ModuleType:
    """Return the factory module for the active project marker.

    Raises if no project marker is set or no factory is registered. The
    pattern mirrors `framework/auth.py` and `framework/api.py`.
    """
    project = _active_project(request)
    if project is None:
        raise RuntimeError(
            "factory fixture requires a project marker (e.g. @pytest.mark.belajar_lintas_ronald)"
        )
    if project not in FACTORIES:
        raise KeyError(f"No factory module registered for project {project!r}")
    return FACTORIES[project]
