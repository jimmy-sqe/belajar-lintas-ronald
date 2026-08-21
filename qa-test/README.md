# qa-test-belajar_lintas_ronald

UI + API automation boilerplate built on pytest + Playwright. Part of the
S-Quantum Engine boilerplate monorepo; scaffolded per product into `qa-test/`.

## Quick start

```bash
python3.12 -m venv .venv   # version pinned by .python-version (3.12)
source .venv/bin/activate
pip install -e ".[dev]"
playwright install chromium
cp .env.example .env        # fill in BASE_URL / TEST_USERNAME / TEST_PASSWORD
pre-commit install
```

## Run tests

```bash
pytest                                               # everything
pytest -m belajar_lintas_ronald                            # one project
pytest tests/belajar_lintas_ronald/test_login.py           # one file
pytest --headed tests/belajar_lintas_ronald/test_login.py  # headed for debugging
```

## Layout

```
framework/        shared primitives (config, auth, api, registries)
  api_clients/    per-project API client subclasses
  factories/      API-first state factories + LIFO cleanup
pages/            page objects, one folder per project
tests/            test code, mirrors pages/
  fixtures/       reusable fixtures (browser, auth, api, factories)
```

## Boilerplate axes

This subproject is pruned from `.boilerplate.yaml`:

| Axis | Options | Effect |
|---|---|---|
| `sample-suite` | `general-service` \| `none` | Ship the sample login suite, or a bare framework. |
| `api-testing` | `wired` \| `none` | UI + API (default), or UI-only. |
| `reporting` | `allure` \| `pytest-html` | pytest-html is always on; `allure` adds Allure. |
| `parallel-exec` | `xdist` \| `none` | pytest-xdist parallel execution, or serial. |

See `CLAUDE.md` for working rules and marker discipline.

## Reports

- **pytest-html** — `pytest-results/report.html` (self-contained).
- **Allure** (if enabled) — `allure-results/`; render with
  `allure serve allure-results`.

Failures attach a screenshot, the URL at failure, and the Playwright trace.
