from __future__ import annotations

from pathlib import Path

from playwright.sync_api import Locator, Page, expect


class BasePage:
    """Optional base class for page objects.

    Subclasses may set PATH (relative; base_url comes from the Playwright
    context) and READY (a Locator that signals the page is interactive).
    """

    PATH: str = ""
    READY: Locator | None = None

    def __init__(self, page: Page):
        self.page = page

    def goto(self) -> None:
        self.page.goto(self.PATH)
        if self.READY is not None:
            expect(self.READY).to_be_visible()

    def screenshot(self, path: str | Path, full_page: bool = True) -> None:
        self.page.screenshot(path=str(path), full_page=full_page)

    def expect_visible(self, locator: Locator, timeout: float | None = None) -> None:
        expect(locator).to_be_visible(timeout=timeout)

    def expect_url_matches(self, pattern: str) -> None:
        expect(self.page).to_have_url(pattern)
