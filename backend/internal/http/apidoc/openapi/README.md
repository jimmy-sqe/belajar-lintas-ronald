# OpenAPI / oapi-codegen

This option uses [oapi-codegen](https://github.com/oapi-codegen/oapi-codegen)
to generate Go types and Echo server stubs from an OpenAPI 3 spec.

## Usage

1. Edit `openapi.yaml` to declare endpoints, request/response shapes, and
   security schemes.
2. Generate code:
   ```
   make oapi-gen
   ```
   This produces `generated.go` next to the spec, containing:
   - Request / response struct types
   - An Echo server interface (`ServerInterface`)
   - Helpers to register the interface against an Echo router

3. Implement the interface in `internal/http/handler/` and pass the impl
   to `RegisterHandlers(e, impl)`.

## When to choose this

- You prefer spec-first / contract-first development.
- Multiple consumers need the contract (e.g., FE, mobile, partners).
- You want type safety on requests / responses without hand-writing DTOs.

## When to choose `swag` instead

- Your spec is small and lives close to code.
- Your team prefers code-first / annotation-driven docs.
