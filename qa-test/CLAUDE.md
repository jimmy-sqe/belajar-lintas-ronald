# qa-test-belajar_lintas_ronald — Working Rules

QA automation boilerplate: pytest + Playwright for UI, requests + faker for
API-first state setup. One product, one sample project (`belajar_lintas_ronald`)
that you rename to your service.

## Layout

```
framework/   shared primitives (config, auth, api, api_clients/, factories/)
pages/       page objects (base + one folder per project)
tests/       test code, mirrors pages/; tests/fixtures/ holds reusable fixtures
```

## Run

```bash
python3.12 -m venv .venv && source .venv/bin/activate
pip install -e ".[dev]"
playwright install chromium
cp .env.example .env          # fill in BASE_URL / TEST_USERNAME / TEST_PASSWORD
pytest                        # everything
pytest -m belajar_lintas_ronald     # one project
pytest --headed tests/belajar_lintas_ronald/test_login.py   # debug one file
```

## Writing a test

Tests live under `tests/<project>/`. The folder name auto-applies the
`@pytest.mark.<project>` marker (see `tests/conftest.py`). Request
`logged_in_page` for a fresh UI login, or `credentials` for raw values.

```python
import pytest
from playwright.sync_api import Page

@pytest.mark.belajar_lintas_ronald
def test_dashboard(logged_in_page: Page):
    logged_in_page.goto("/dashboard")
```

## API-first state setup

Arrange preconditions through the `api` + `factory` + `cleanup_registry`
fixtures, not the UI. `cleanup_registry` drains LIFO at teardown.

## Conventions

- **Reruns are opt-in.** Tag a genuinely flaky test `@pytest.mark.flaky`
  (≤ 2 reruns). Don't tag a test just because it failed once.
- **Smoke set ≤ 10 tests/project.** `@pytest.mark.smoke` marks critical-path
  tests answering "is this broken right now?" in under a minute.
- **Reports:** pytest-html → `pytest-results/report.html`; Allure (if enabled)
  → `allure-results/` (`allure serve allure-results`).
- **Lint/type:** `ruff check .` and `mypy framework pages tests` must pass
  (mypy is strict). `pre-commit install` wires both.

## Boilerplate axes (marker discipline)

This is a pruned boilerplate. Wiring files carry axis markers:

```python
# boilerplate:axis=api-testing option=wired START
from framework.api import api_for
# boilerplate:axis=api-testing option=wired END
```

When importing axis-scoped content into a kept file, wrap **both** the import
AND every usage. Axes: `sample-suite`, `api-testing`, `reporting`,
`parallel-exec` — see `.boilerplate.yaml`. Keep marker pairs balanced; an
unbalanced pair breaks the pruner.

## Adding a project

1. `pages/<slug>/login_page.py` subclassing `BaseLoginPage`; register in
   `LOGIN_PAGES` (`framework/auth.py`).
2. API client (when needed): `framework/api_clients/<slug>.py` subclassing
   `BaseAPIClient`; register in `API_CLIENTS` (`framework/api.py`).
3. Add `<slug>_base_url` / `<slug>_test_username` / `<slug>_test_password` /
   `<slug>_viewport` to `Settings` (`framework/config.py`) and `<slug>` to
   `PROJECTS`.
4. Add a `<slug>: ...` marker to `pyproject.toml [tool.pytest.ini_options]`.
