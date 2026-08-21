"""Factory + cleanup primitives for API-first state setup.

Convention:
    * Tests that need state should arrange it through `factory(...)` + `api`,
      not by driving the UI. The UI test should start at the state under test.
    * Each `create_*` call registers a cleanup with the active `CleanupRegistry`
      so teardown is automatic.
"""

from __future__ import annotations

from collections.abc import Callable
from dataclasses import dataclass, field
from typing import Any, Protocol

from framework.api_clients.base import BaseAPIClient


@dataclass
class CleanupRegistry:
    """LIFO registry of (label, callable) cleanup actions.

    Callers register via `add(...)`; the test fixture drains via `drain()` on
    teardown. Cleanup failures are swallowed but logged to the registry so a
    test can still pass if e.g. an entity was already deleted as part of the
    test body.
    """

    _actions: list[tuple[str, Callable[[], Any]]] = field(default_factory=list)
    errors: list[tuple[str, BaseException]] = field(default_factory=list)

    def add(self, label: str, fn: Callable[[], Any]) -> None:
        self._actions.append((label, fn))

    def drain(self) -> None:
        while self._actions:
            label, fn = self._actions.pop()
            try:
                fn()
            except BaseException as exc:  # noqa: BLE001 — teardown must not raise
                self.errors.append((label, exc))


class Factory(Protocol):
    """A factory takes an API client + overrides and returns the created resource.

    The factory MUST register a cleanup action on the supplied CleanupRegistry
    before returning.
    """

    def __call__(
        self,
        api_client: BaseAPIClient,
        cleanup: CleanupRegistry,
        **overrides: Any,
    ) -> dict[str, Any]: ...
