# Problem Model

## Responsibility

The problem module keeps problem identity separate from optional content and
from the network operation that fetches content.

```go
type Problem struct {
	ID      string
	Content *string
}
```

`ID` is the only field required for source-code submission. `Content == nil`
means content is absent or was not downloaded; absence never blocks
submission.

## Data-only model

`Problem` does not:

- fetch itself;
- tokenize or normalize its ID;
- parse HTML;
- enforce language restrictions;
- submit source code.

The CLI creates `Problem{ID: pid, Content: nil}` directly for `submit`. This is
valid because the verified direct-submission flow does not request
`GET /problem/{pid}`.

## Optional loader

Problem-content fetching is a public plugin port because HTML extraction is
optional and can evolve independently:

```go
type ProblemLoader interface {
	Load(ctx context.Context, problemID string) (Problem, error)
}
```

The loader implementation uses an authorized transport capability supplied by
its factory. It returns a new `Problem` and does not mutate an existing value.

The first-release CLI never selects or invokes a problem loader. It exists so
library consumers can request content without coupling submission to a page
request.

## Content semantics

- A successful content load returns a non-nil `Content`.
- A page that is valid but has no extracted statement returns a typed parse or
  content error; it does not silently return empty content as success.
- Language hints extracted from a problem page are advisory. The runtime
  catalog remains authoritative for whether the client may submit.
- Content is treated as untrusted text and is not interpolated into logs.
- The loader bounds response size and parsing time.

## Authentication

The problem endpoint requires an authenticated session. The selected loader
receives only an authorized request capability; it cannot read a stored
password or replace the session's cookie jar.

## Tests

- `Problem{ID: "<pid>", Content: nil}` is a valid submission input.
- Submission does not call the loader.
- Loader errors never mutate a supplied problem.
- A malformed page cannot be represented as a successful empty load.
- Advisory page language restrictions cannot override runtime language policy.
- Unauthorized responses map to `ErrUnauthenticated`.
