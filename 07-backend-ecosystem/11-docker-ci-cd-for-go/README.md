# Docker & CI/CD for Go

[← Section index](../README.md) · Prev: [Security: mTLS, JWT, OAuth & Istio](../10-security-mtls-jwt-oauth-istio/README.md)

**Interview weight: 0.8%**  ·  Revise in: 5 min  ·  Run: `go run ./07-backend-ecosystem/11-docker-ci-cd-for-go`

## In one line
Build the Go binary in one Docker stage, copy only the binary into a tiny final image, and let CI run format, vet, test and build on every push.

## Must know
- Multi-stage build: a `golang` image compiles, a small runtime image (`distroless` or `scratch`) runs. The final image holds just the binary.
- `CGO_ENABLED=0` gives a static binary that runs on `scratch` or distroless with no libc.
- Copy `go.mod` and `go.sum` first and run `go mod download`, so Docker caches dependencies until they change.
- Run as a non-root user. Add a `.dockerignore` (`.git`, test data, local files).
- Config via env vars (12-factor), not baked into the image. Secrets come from the platform, never from the repo or image layers.
- Inject the version at build time: `-ldflags "-X main.version=1.2.3"`. `-s -w` strips debug info to shrink the binary.
- Tag images with the git SHA or version, not only `latest`.
- A typical CI pipeline: `gofmt` check, `go vet`, `go test -race ./...`, lint (`golangci-lint`), build, build and push image, deploy.

## Code to remember
```dockerfile
FROM golang:1.25 AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /app .

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /app /app
USER nonroot
ENTRYPOINT ["/app"]
```

## Common traps
- Copying the whole source before `go mod download`, so every code change re-downloads dependencies.
- A dynamically linked binary on `scratch` fails to start with "not found".
- Running as root, or shipping the Go toolchain and source in the final image.
- Using `latest` in deployments, so you cannot tell what is running or roll back.

## Interview Q&A

**Q1. Why multi-stage builds?**

**A:** Smaller and safer final images (tens of MB instead of hundreds), no compiler or source code in production.

**Q2. Why `CGO_ENABLED=0`?**

**A:** It produces a static binary with no dependency on system C libraries, so it runs in minimal images.

**Q3. What does your CI pipeline do?**

**A:** Check formatting, vet, run tests with `-race`, lint, build the image, push it with a unique tag, then deploy. Fail fast on the cheap checks first.

**Q4. CI vs CD, and how do you roll back?**

**A:** CI builds and tests every push. CD ships the passing build to an environment. Roll back by redeploying the previous image tag, which is why tags are unique.

**Q5. How do you pass configuration to the container?**

**A:** Environment variables, ConfigMaps or Secrets in Kubernetes. The same image runs in every environment.
