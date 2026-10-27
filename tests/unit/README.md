# Unit Tests

The current test tree is intentionally incomplete. It establishes locations
and a shared suite harness without asserting product behavior.

```text
tests/unit/
  functional/   functional unit-test suite
  regression/   regression unit-test suite
```

Coverage is measured across the whole module while running both unit-test
categories. It is a test concern, not a separate product test type.

```sh
go test ./tests/unit/...
go test ./tests/unit/functional
go test ./tests/unit/regression
go test ./tests/unit/... -coverpkg=./... -covermode=atomic -coverprofile=coverage.out
go tool cover -func=coverage.out
```

The initial suites contain no cases. Tests for integrations, end-to-end
workflows, and live upstream services are not part of this scaffold.
