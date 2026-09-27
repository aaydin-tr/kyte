# Kyte

A Go library for composing MongoDB Queries through a chainable Filter. Kyte only produces Queries; it never talks to a MongoDB server itself.

## Language

**Filter**:
The chainable builder returned by `kyte.Filter()`, on which operators are recorded.
_Avoid_: Query builder, builder

**Query**:
The MongoDB query filter document (`bson.D`) that a Filter's `Build()` returns.
_Avoid_: Filter (for the built document), filter document

**Source**:
The struct pointer a Filter is given so that pointers to its fields resolve to bson keys.
_Avoid_: Schema, model

**Global Filter**:
A Filter that is automatically merged into every new Filter unless that Filter opts out.
_Avoid_: Default filter, tenant filter
