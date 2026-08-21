from __future__ import annotations

from pathlib import Path

from pydantic import BaseModel
from pydantic_settings import BaseSettings, SettingsConfigDict

_REPO_ROOT = Path(__file__).resolve().parent.parent

# Single-product boilerplate: one sample project ships by default. The
# scaffolder renames the sample slug to the product slug. When the
# sample-suite axis is pruned, this tuple becomes empty and the framework
# falls back to the global BASE_URL / TEST_* settings.
PROJECTS: tuple[str, ...] = (
    # boilerplate:axis=sample-suite option=general-service START
    "belajar_lintas_ronald",
    # boilerplate:axis=sample-suite option=general-service END
)


class ProjectConfig(BaseModel):
    base_url: str
    username: str
    password: str
    viewport: tuple[int, int]


def _parse_viewport(value: str | None) -> tuple[int, int] | None:
    if not value:
        return None
    try:
        w, h = value.lower().split("x", 1)
        return (int(w.strip()), int(h.strip()))
    except (ValueError, AttributeError):
        return None


class Settings(BaseSettings):
    # env_file is absolute so .env loads regardless of pytest's working dir
    # (e.g. when an IDE runs a single test file from tests/<project>/).
    model_config = SettingsConfigDict(
        env_file=_REPO_ROOT / ".env", extra="ignore", case_sensitive=False
    )

    base_url: str = ""
    test_username: str = ""
    test_password: str = ""
    viewport: str = "1440x900"

    # boilerplate:axis=sample-suite option=general-service START
    belajar_lintas_ronald_base_url: str | None = None
    belajar_lintas_ronald_test_username: str | None = None
    belajar_lintas_ronald_test_password: str | None = None
    belajar_lintas_ronald_viewport: str | None = None
    # boilerplate:axis=sample-suite option=general-service END

    def project(self, name: str) -> ProjectConfig:
        if name not in PROJECTS:
            raise KeyError(f"Unknown project: {name!r}. Known: {PROJECTS}")

        def get(k: str) -> str | None:
            return getattr(self, f"{name}_{k}", None)

        viewport = _parse_viewport(get("viewport")) or _parse_viewport(self.viewport) or (1440, 900)
        return ProjectConfig(
            base_url=get("base_url") or self.base_url,
            username=get("test_username") or self.test_username,
            password=get("test_password") or self.test_password,
            viewport=viewport,
        )
