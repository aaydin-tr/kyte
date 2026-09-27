# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project

Kyte (`github.com/aaydin-tr/kyte`) is a single-package Go library: a fluent MongoDB filter query builder on top of `go.mongodb.org/mongo-driver` (v1, `bson` package). Only filter operations exist; aggregate/update builders are planned but not implemented. Go 1.20 (CI uses 1.20.1), so avoid newer language/stdlib features.

## Commands

```sh
go build -v ./...                                   # build (what CI runs)
go test -v ./...                                    # all tests (what CI runs)
go test -run TestFilter_Equal ./...                 # single test
go test -run 'TestFilter_Equal/without_source' ./...  # single subtest
go test -bench . -run '^$' ./...                    # benchmarks (filter_benchmark_test.go)
go test -coverprofile=coverage.out ./...            # coverage
```

Releases are tag-driven via GoReleaser (`builds: skip: true` — library only, no binaries). Changelog excludes commits prefixed `docs:`, `test:`, `chore:`; the repo uses conventional-commit prefixes (`feat:`, `fix:`, `docs:`).

## Architecture

Two source files:

- [kyte.go](kyte.go) — the `kyte` struct: schema/field resolution and validation. `setSourceAndPrepareFields` reflects over the `Source` struct pointer and builds `fields map[any]string` mapping **the address of each struct field** (e.g. `&user.Name`) to its bson key, recursing into nested structs, pointers-to-struct and slices of structs to produce dotted paths (`parent.child`). Bson tag parsing skips builtin flags (`omitempty`, `minsize`, `truncate`, `inline`); fields tagged `-` or untagged are ignored. All exported errors (`Err*`) live here. Options use the functional-option pattern (`Source`, `ValidateField`, `IgnoreGlobalFilters`).
- [filter.go](filter.go) — the public `filter` builder, operator methods, and global filters.

Key behaviors that span both files:

- **Field arguments are `any`**: either a string key or a pointer to a field of the `Source` struct. Pointers are resolved to bson keys via the `fields` map; without a `Source`, only strings work (`ErrFieldMustBeString`). Validation against the schema is on by default when a `Source` is given.
- **Deferred vs. immediate**: most operator methods just record an `operation` via `f.set(...)`; validation and conversion to `bson.D` happen in `Build()`. But global filters (applied inside `Filter()`), `Raw`, and `And`/`Or`/`NOR` append to `f.query` immediately — so they appear before deferred operations in the output.
- **Errors are sticky**: `kyte.setError` keeps only the first error; methods keep returning `f` for chaining and `Build()` returns that error. New operators should follow this pattern rather than panicking or returning errors directly.
- **`Build()` is memoized** via `isBuild`; later calls return the cached query.
- **Nested logical filters** (`And`/`Or`/`NOR`) re-point the child filter's `kyte` at the parent's `Source` and `checkField`, build it, and wrap each resulting element as `bson.M` in a `bson.A`. Child filters created with `Filter()` also receive global filters unless `IgnoreGlobalFilters()` is passed.
- **Value normalization in `Build()`**: pointer values are dereferenced; non-slice values for `$in`/`$nin`/`$type` are wrapped in `bson.A`; `$regex` stores its full `{$regex, $options}` doc as the field value; `$where` and `$jsonSchema` are top-level and need no field.
- **Global filters** (`AddGlobalFilter`, `ClearGlobalFilters`, `GetGlobalFilters`) are package-level state guarded by `globalMutex`.

Operator constants are at the top of [filter.go](filter.go), including TODO lists of unimplemented operators (`$elemMatch`, `$not`, `$text`, `$expr`, geo, bitwise).

## Tests

- [filter_test.go](filter_test.go) is an external test package (`kyte_test`) exercising the public API; [kyte_test.go](kyte_test.go) and [filter_benchmark_test.go](filter_benchmark_test.go) are internal (`package kyte`) and access unexported types.
- Tests use plain `testing` (no assertion library) with `t.Parallel()` and `t.Run` subtests. Because global filters are shared package state, any test that adds global filters can leak into parallel tests — clear them and avoid `t.Parallel()` in such tests.

## Docs

[README.md](README.md) documents every operator with its generated query; update it when adding or changing operators. Note the README uses `Nor`/`JsonSchema` while the actual method names are `NOR`/`JSONSchema`.

## Agent skills

### Issue tracker

Issues live as local markdown files under `.scratch/<feature>/` (gitignored, local-only). See `docs/agents/issue-tracker.md`.

### Triage labels

Default five-role vocabulary (`needs-triage`, `needs-info`, `ready-for-agent`, `ready-for-human`, `wontfix`). See `docs/agents/triage-labels.md`.

### Domain docs

Single-context: one root `CONTEXT.md` + `docs/adr/`. See `docs/agents/domain.md`.
