# ADR 0002: Critical Invariants Are Enforced in the Database

## Status
Accepted

## Decision
Critical borrow and stock invariants are enforced in PostgreSQL with constraints and indexes, not only in application code.

## Why
- Application checks alone are not enough under concurrent requests.
- Database constraints prevent invalid state from being persisted.
- Partial unique index and stock constraints act as the last safety net.
