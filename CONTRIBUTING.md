# Contributing to GoServe

GoServe is intentionally a small foundation repository. Contributions are welcome, but the bar is different from an application repo:

- improve the foundation
- keep the starter minimal
- prefer clarity and testability over cleverness
- avoid adding product-specific behavior to the base project

## Good First Steps

1. Read [README.md](/Users/admin/Projects/opensource/org/goserve/README.md).
2. Read [docs/good-first-issues.md](/Users/admin/Projects/opensource/org/goserve/docs/good-first-issues.md).
3. Set up local keys with `make keys`.
4. Run validation with `go test ./...` and `go build ./...`.
5. Pick one focused change.

## Development Workflow

1. Create a branch from `main`.
2. Make a small, well-scoped change.
3. Add or update tests when behavior changes.
4. Run:

```bash
go test ./...
go build ./...
```

5. Update documentation if the change affects contributor or runtime behavior.
6. Open a pull request with a clear summary and verification notes.

## What Belongs Here

Good contributions:

- small improvements to auth, middleware, logging, and startup behavior
- deterministic starter routes and tooling improvements
- tests around existing foundation behavior
- contributor documentation and developer experience improvements
- infrastructure fixes that make local or CI validation easier

Usually out of scope:

- product/domain features
- large dependency additions without strong justification
- turning the starter into a full application
- adding demo endpoints that do not exercise reusable foundation patterns

## Coding Expectations

- Keep changes small and readable.
- Preserve the existing layered structure: `foundation`, `business`, `app`.
- Use ASCII unless a file already requires something else.
- Prefer deterministic behavior over randomness in starter code.
- Add comments only when they clarify non-obvious logic.
- Keep public docs aligned with the current implementation.

## Tests

At minimum, run:

```bash
go test ./...
go build ./...
```

If you touch auth or HTTP behavior, include focused tests near the changed package.

## Pull Requests

A good pull request includes:

- what changed
- why it changed
- how it was verified
- any scope decisions or tradeoffs

Small PRs are much easier to review and merge than broad refactors.

## Issue Selection

If you are new to the repo, start with one item from [docs/good-first-issues.md](/Users/admin/Projects/opensource/org/goserve/docs/good-first-issues.md). Those are intentionally scoped so contributors can learn the codebase without pushing product logic into the starter.
