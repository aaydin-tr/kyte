# Integration tests live in the root module; Kyte requires Go 1.25

Integration tests execute Queries against a real MongoDB started by testcontainers-go, whose current releases require Go 1.25. We keep those tests in the root module and raised Kyte's minimum Go from 1.20 to 1.25, accepting that users on older Go can no longer upgrade Kyte, in exchange for a single module and a single `go.mod` to maintain.

## Considered Options

- **Nested `integration/` module on Go 1.25** — would have kept the library on Go 1.20, at the cost of a second module, a `replace` directive, and two Go versions in CI.
- **GitHub Actions `services: mongo` plus a connection-string env var** — no new dependencies, but local runs need a hand-started container and the MongoDB setup lives in CI YAML rather than in the tests.

## Consequences

- testcontainers-go is a requirement in Kyte's `go.mod`, so it appears in consumers' module graphs even though only tests import it.
- Only testcontainers-go core is used, with a generic `mongo` container. Its `modules/mongodb` helper was rejected because it requires mongo-driver v2, whose client cannot take Kyte's v1 `bson.D` as a filter, and would add a second driver major to every consumer's graph.
