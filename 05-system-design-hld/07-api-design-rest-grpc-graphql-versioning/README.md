# API Design: REST vs gRPC vs GraphQL & Versioning

[← Section index](../README.md) · Prev: [Microservices Communication & Resilience](../06-microservices-communication-and-resilience/README.md) · Next: [Consistent Hashing & Distributed Cache](../08-consistent-hashing-and-distributed-cache/README.md)

**Interview weight: 1.1%**  ·  Revise in: 5 min

## In one line
Use REST for public APIs, gRPC for internal service calls, GraphQL when clients need flexible queries, and version carefully so old clients keep working.

## Must know
- REST: stateless, resource URLs (`/orders/42`), HTTP verbs, proper status codes (201 created, 400 bad input, 401 not logged in, 403 not allowed, 404, 409 conflict, 429, 5xx).
- Idempotent (repeating gives the same result) verbs: GET, PUT, DELETE. POST is not, so use an `Idempotency-Key` for payments and orders.
- gRPC: protobuf over HTTP/2, typed contract in `.proto`, four call types (unary, server-stream, client-stream, bidirectional), deadlines built in.
- GraphQL: one endpoint, the client picks fields. Watch for N+1 queries (fix with batching, DataLoader style) and expensive queries (depth or cost limits).
- Pagination: cursor-based beats offset for large or changing data. Offset gets slow and can skip or repeat rows.
- Versioning: URL prefix `/v1/` is simplest. A header is cleaner but harder to test. Prefer additive changes (new optional fields) so you rarely need a new version.
- Deprecate with a notice and a `Sunset` header, keep the old version alive during migration.
- Protobuf compatibility: never reuse or renumber a field; add new fields, reserve removed ones.

## Common traps
- Using POST for everything, or returning 200 with an error body.
- Offset pagination on a table that changes while the user pages.
- Breaking changes (renaming or removing a field) without a version bump.
- Exposing gRPC directly to browsers without a gateway.

## Interview Q&A

**Q1. When would you choose gRPC over REST?**

**A:** Internal service-to-service calls that need low latency, streaming, or a strict contract. REST is better for public APIs and browser clients.

**Q2. What is the N+1 problem in GraphQL?**

**A:** Resolving a list of N items where each item triggers its own DB query, so N+1 queries. Fix by batching lookups per request (DataLoader pattern).

**Q3. How do you change a REST API without breaking clients?**

**A:** Make additive changes where possible. For breaking ones, add a new version, mark the old one deprecated with a `Sunset` header, and keep it running during the transition.

**Q4. Cursor vs offset pagination?**

**A:** Cursor uses the last seen key (`WHERE id > ? ORDER BY id LIMIT n`), stays fast and stable as data changes. Offset scans and discards rows, and shifts when rows are inserted.

**Q5. POST vs PUT vs PATCH: which are idempotent?**

**A:** POST creates or triggers an action, not idempotent (add an idempotency key to make retries safe). PUT replaces the whole resource, idempotent. PATCH changes some fields.
