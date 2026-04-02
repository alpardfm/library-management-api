# Release Prep

## Checklist

1. Ensure `make quality` passes locally.
2. Run optional integration test if PostgreSQL test DB is available.
3. Review `CHANGELOG.md` and move relevant items from `Unreleased` into a versioned section.
4. Confirm `.env.example`, README, and workflow docs still match the project behavior.
5. Create and push tag:

```bash
git tag -a v0.1.0 -m "Release v0.1.0"
git push origin v0.1.0
```

## Notes

- PR quality checks are provided by `.github/workflows/quality-gates.yml`.
- Optional integration CI can be triggered manually from GitHub Actions.
