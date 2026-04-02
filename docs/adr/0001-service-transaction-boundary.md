# ADR 0001: Service Owns Transaction Boundary

## Status
Accepted

## Decision
Transaction boundaries are created in the service layer. Repository methods are transaction-aware via `WithTx(...)`, but business flow remains in services.

## Why
- Borrow and return touch multiple repositories.
- Business validation should stay outside the repository layer.
- The service layer is the right place to keep one consistency boundary for one use case.
