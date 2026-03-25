# GoServe

GoServe is a shareable Go service starter for building authenticated HTTP APIs with a small, explicit architecture. It is intentionally a starter, not a finished product: the repo gives you the service skeleton, middleware stack, auth wiring, local key workflow, and deployment base so new feature work can happen outside the foundation.

## What Exists Today

- Layered project structure across `foundation`, `business`, and `app`
- HTTP service with health endpoints, debug endpoints, and starter example routes
- JWT authentication backed by RSA keys and OPA policy evaluation
- Structured logging, metrics, panic recovery, and graceful shutdown
- Docker build and Kubernetes manifests
- Local admin CLI for key generation and JWT creation
- Automated tests for the starter auth and handler flows

## Project Layout

```text
goserve/
├── app/
│   ├── services/sales-api/      # HTTP service entrypoint
│   └── tooling/sales-admin/     # Local key + token helper
├── business/
│   └── web/v1/                  # Auth, middleware, HTTP groups
├── foundation/
│   ├── keystore/                # PEM-backed key loading + generation
│   ├── logger/                  # Structured logging
│   └── web/                     # Lightweight web primitives
└── zarf/
    ├── docker/                  # Container build
    └── k8s/                     # Base and dev manifests
```

## Quick Start

Prerequisites:

- Go 1.25.5+
- Make
- Docker and Kubernetes tooling if you want the container or K8s flows

Generate a local dev key:

```bash
make keys
```

Run the service:

```bash
make run
```

Hit the public starter route:

```bash
make curl
```

Generate a token and call the authenticated route:

```bash
TOKEN="$(make token)"
make curl-auth TOKEN="$TOKEN"
```

## Starter Routes

- `GET /v1/readiness`
- `GET /v1/liveness`
- `GET /v1/example`
- `GET /v1/exampleauth`
- `GET /debug/vars`
- `GET /debug/pprof/*`

`/v1/exampleauth` requires a valid bearer token with the `ADMIN` role.

## Local Auth Workflow

The repo no longer ships private keys. Local development keys are generated into `.local/keys/`, which is gitignored.

Generate a key:

```bash
go run app/tooling/sales-admin/main.go -action generate-keys
```

Generate a token from the first local key:

```bash
go run app/tooling/sales-admin/main.go -action generate-token
```

You can override the folder, key id, issuer, subject, roles, and TTL with flags.

## Validation

Run the full local verification pass:

```bash
go test ./...
go build ./...
```

## Deployment Notes

- The runtime image does not embed auth keys.
- Production keys should be mounted or injected by your deployment platform.
- The Kubernetes manifests are a base to extend, not a claim of production completeness.

## Next Step for Product Teams

This repository should stay focused on foundation concerns. New domain behavior should be added as dedicated route groups, business packages, and tests rather than as placeholder endpoints inside the starter.
