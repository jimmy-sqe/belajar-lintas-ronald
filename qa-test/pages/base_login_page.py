from abc import ABC, abstractmethod

from playwright.sync_api import Page


class BaseLoginPage(ABC):
    """Each project subclasses this and implements the full login flow.

    The contract is intentionally one method: projects with 2-step flows,
    SSO redirects, magic links, or MFA cannot fit a fixed multi-step API.
    """

    def __init__(self, page: Page):
        self.page = page

    @abstractmethod
    def perform_login(self, username: str, password: str) -> None:
        """Run the full login flow; return only when post-login state is visible."""
