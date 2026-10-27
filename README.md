# LGC

LGC is a Luogu code submission project designed to let an LLM submit code to
Luogu and inspect the judging result. It targets workflows where code boards or
solution snippets need to be tested against real problems rather than only
local unit tests.

The intended workflow is:

1. Authenticate with the Luogu web service through an explicitly configured
   CAPTCHA solver.
2. Resolve the problem identity and runtime submission state; opening the
   problem page is optional for source-code submission.
3. Resolve the correct language ID from Luogu's runtime configuration.
4. Submit source code and receive a judge record ID.
5. Wait for asynchronous judging through record polling or WebSocket updates.
6. Return a machine-readable verdict, diagnostics, and test progress to the
   caller.

## Project status

The repository now contains the Go project skeleton. Public library values,
provider ports, plugin entry points, and the four-stage CLI pipeline are
declared, but their behavior is intentionally unimplemented. Entry points that
require business logic return `domain.ErrNotImplemented`.

The test tree currently contains an empty unit-test suite only. Functional and
regression suites are separate directories, and coverage uses the standard Go
test tooling over both suites. Integration and live upstream tests are not yet
scaffolded.

All project documentation and contribution requirements are maintained under
[`docs/`](docs/).

The proposed client architecture starts at
[`docs/design.md`](docs/design.md), with component and library-module designs
under [`docs/design/`](docs/design/).

The initial CLI design exposes only authentication and submission. It does not
expose problem-content fetching.

## Repository layout

```text
cmd/       CLI entry point and four-stage command pipeline
lib/       domain values, provider ports, and core workflow declarations
plugin/    bundled provider entry-point skeletons
tests/     unit-test suites; see tests/unit/README.md
docs/      design and upstream API documentation
```

## Development

The project targets Go 1.22 or newer.

```sh
go test ./...
go test ./tests/unit/... -coverprofile=coverage.out
go tool cover -func=coverage.out
```

The executable package is `./cmd/lgc`. It is not yet behaviorally usable while
the command pipeline remains a declaration-only scaffold.
