# swaggo / swag annotations

This option uses [swaggo/swag](https://github.com/swaggo/swag) to generate
Swagger UI from inline Go comments on handlers.

## Usage

1. Annotate handlers with swag tags (see `example.go` for the pattern).
2. Generate docs:
   ```
   make swag-gen
   ```
3. The generated docs land at `docs/swagger.json` / `docs/swagger.yaml`
   (configure via `swag init -o`).
4. Serve them via [echo-swagger](https://github.com/swaggo/echo-swagger)
   in `internal/http/router.go` if desired.

## When to choose this

- Handlers are short and benefit from inline doc comments.
- Spec is generated from code (code-first workflow).

## When to choose `oapi-codegen` instead

- You prefer spec-first development.
- You want to generate types and server stubs from `openapi.yaml`.
