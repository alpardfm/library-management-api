# ADR 0003: API Uses a Standard Response Envelope

## Status
Accepted

## Decision
HTTP handlers return a shared response envelope with `success`, `message`, `data`, `error`, and `meta`.

## Why
- Clients get one predictable shape across endpoints.
- Error-to-status mapping stays consistent.
- Pagination and list metadata can be added without changing endpoint-specific contracts.
