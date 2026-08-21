from __future__ import annotations

from playwright.sync_api import Page

from framework.config import Settings
from pages.base_login_page import BaseLoginPage
# boilerplate:axis=sample-suite option=general-service START
from pages.belajar_lintas_ronald.login_page import LoginPage
# boilerplate:axis=sample-suite option=general-service END

LOGIN_PAGES: dict[str, type[BaseLoginPage]] = {
    # boilerplate:axis=sample-suite option=general-service START
    "belajar_lintas_ronald": LoginPage,
    # boilerplate:axis=sample-suite option=general-service END
}


def login_as(page: Page, project: str, settings: Settings) -> None:
    """Perform a fresh UI login for the given project marker."""
    if project not in LOGIN_PAGES:
        raise KeyError(f"No LoginPage registered for project {project!r}")
    cfg = settings.project(project)
    LOGIN_PAGES[project](page).perform_login(cfg.username, cfg.password)
