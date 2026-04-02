# Changelog

All notable changes to this project will be documented in this file.

## Unreleased

### Added
- Request ID middleware and structured request logging.
- Standard response envelope and domain error mapping.
- Query parsing helper with consistent pagination metadata.
- Transaction-aware borrow and return flow.
- Database invariants and supporting indexes for stock and active borrows.
- GitHub Actions quality gate for lint and unit tests.

### Changed
- Return policy is now role-aware for `admin`, `librarian`, and `member`.
- Integration and E2E test setup now skips cleanly when environment is unavailable.
- README, Makefile, and CI docs updated for faster onboarding.
