# Good First Issues

These issues are intentionally scoped for new contributors. They improve the starter without turning it into an application.

## 1. Add `make test` and `make check`

Goal:
- add a simple contributor-friendly test target
- add a `check` target that runs the basic validation flow

Suggested scope:
- update [makefile](/Users/admin/Projects/opensource/org/goserve/makefile)
- document the commands in [README.md](/Users/admin/Projects/opensource/org/goserve/README.md) and [CONTRIBUTING.md](/Users/admin/Projects/opensource/org/goserve/CONTRIBUTING.md)

Acceptance criteria:
- `make test` runs `go test ./...`
- `make check` runs at least tests and build
- docs mention the new targets

Labels:
- `good first issue`
- `documentation`
- `developer experience`

## 2. Add tests for health endpoints

Goal:
- add focused tests for readiness and liveness handlers

Suggested scope:
- add tests near [app/services/sales-api/v1/handlers/checkgrp](/Users/admin/Projects/opensource/org/goserve/app/services/sales-api/v1/handlers/checkgrp)

Acceptance criteria:
- readiness test verifies `200` and response body
- liveness test verifies `200` and basic fields

Labels:
- `good first issue`
- `testing`

## 3. Add unauthorized and forbidden middleware tests

Goal:
- expand HTTP auth coverage for middleware behavior

Suggested scope:
- add tests near [business/web/v1/mid](/Users/admin/Projects/opensource/org/goserve/business/web/v1/mid)

Acceptance criteria:
- test missing bearer token returns `401`
- test failed authorization returns `403`

Labels:
- `good first issue`
- `testing`
- `auth`

## 4. Add request ID examples to the README

Goal:
- make observability features easier for contributors to understand

Suggested scope:
- update [README.md](/Users/admin/Projects/opensource/org/goserve/README.md)
- optionally include one short sample log line

Acceptance criteria:
- README explains trace IDs and where they appear
- docs stay accurate to the current implementation

Labels:
- `good first issue`
- `documentation`

## 5. Add a PR template

Goal:
- improve contribution quality without adding process overhead

Suggested scope:
- add `.github/pull_request_template.md`

Acceptance criteria:
- template asks for summary, verification, and scope notes
- template stays short and practical

Labels:
- `good first issue`
- `documentation`
- `maintenance`

## 6. Add CLI usage examples for `sales-admin`

Goal:
- document the local auth workflow more clearly

Suggested scope:
- update [README.md](/Users/admin/Projects/opensource/org/goserve/README.md)
- update [CONTRIBUTING.md](/Users/admin/Projects/opensource/org/goserve/CONTRIBUTING.md)

Acceptance criteria:
- docs show key generation and token generation examples
- docs mention how to override `kid`, `issuer`, `subject`, and `ttl`

Labels:
- `good first issue`
- `documentation`
- `tooling`

## 7. Add a small test helper for authenticated HTTP requests

Goal:
- reduce duplicated auth setup in future HTTP tests

Suggested scope:
- add a tiny helper in the relevant test package
- keep it local and minimal

Acceptance criteria:
- helper is used by at least one existing or new test
- helper does not introduce a broad test framework abstraction

Labels:
- `good first issue`
- `testing`

## 8. Add a repository architecture diagram section

Goal:
- help new contributors understand layering faster

Suggested scope:
- update [README.md](/Users/admin/Projects/opensource/org/goserve/README.md)

Acceptance criteria:
- explain `foundation -> business -> app`
- include a short “where should this change go?” guide

Labels:
- `good first issue`
- `documentation`

## Maintainer Note

When opening one of these as a real GitHub issue:

- keep the scope narrow
- include acceptance criteria
- attach the `good first issue` label
- avoid assigning work that adds product/domain logic to the starter
