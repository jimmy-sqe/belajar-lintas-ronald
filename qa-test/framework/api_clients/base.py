from abc import ABC, abstractmethod
from typing import Any

import requests


class BaseAPIClient(ABC):
    """Each project subclasses this and implements its login flow.

    Subclasses may also expose domain helpers (create_finding, delete_user, etc.).
    `self.s` is a requests.Session that carries auth state after login() returns.
    """

    def __init__(self, base_url: str, session: requests.Session | None = None) -> None:
        self.base_url = base_url.rstrip("/")
        self.s = session if session is not None else requests.Session()

    @abstractmethod
    def login(self, username: str, password: str) -> None:
        """Authenticate; on return self.s carries whatever auth state the API needs."""

    def _url(self, path: str) -> str:
        return f"{self.base_url}{path}"

    def get(self, path: str, **kw: Any) -> requests.Response:
        return self.s.get(self._url(path), **kw)

    def post(self, path: str, **kw: Any) -> requests.Response:
        return self.s.post(self._url(path), **kw)

    def put(self, path: str, **kw: Any) -> requests.Response:
        return self.s.put(self._url(path), **kw)

    def delete(self, path: str, **kw: Any) -> requests.Response:
        return self.s.delete(self._url(path), **kw)
