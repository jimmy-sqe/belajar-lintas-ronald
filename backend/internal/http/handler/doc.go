// Package handler holds Echo HTTP handlers. Handlers bind requests, validate
// DTOs, call domain services, and map errors to HTTP responses. They MUST
// NOT import from internal/adapter directly — only from internal/domain.
package handler
