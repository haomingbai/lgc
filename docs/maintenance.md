# Project Maintenance

This document defines the maintenance requirements for the Luogu submission
client's documentation, source comments, API reference, and commit history.

## Language

- Write project documentation in English.
- Write source-code comments in English.
- Keep headings, prose, tables, examples, and explanatory comments in English.
- Preserve protocol identifiers, endpoint paths, field names, and quoted
  upstream messages exactly when they are part of an observed contract.
- Explain a non-English upstream value in English when its meaning affects
  client behavior.

Comments should explain why a constraint exists or why a non-obvious approach
is required. Do not use comments to restate code.

```go
// refreshChallenge obtains new challenge state because CAPTCHA answers and
// C3VK values are single-use.
func refreshChallenge() {
	// Fetch a new challenge from the upstream service.
}
```

Go comments should appear immediately above the declaration they describe.
Documentation comments for exported identifiers should begin with the
identifier's name. Keep comments focused on intent, constraints, and
non-obvious behavior.

## Documentation

- Keep the API reference and this maintenance document current.
- Update the relevant reference in the same change as the behavior it
  describes.
- Use relative links for repository files.
- Use fenced code blocks with a language identifier where possible.
- Remove credentials, session cookies, CSRF tokens, CAPTCHA answers, personal
  identifiers, and record IDs from examples and screenshots.
- Keep one authoritative explanation and link to it rather than duplicating
  detailed behavior across several documents.

## API Reference Maintenance

The upstream service is not a stable public API. Treat every documented
contract as time-sensitive.

When changing an endpoint contract:

1. Capture or identify reproducible evidence for the request and response.
2. Record the observation date, protocol, request method, path, relevant
   status codes, and headers or body fields that affect the contract.
3. Separate verified behavior from inference and unverified behavior.
4. Update cookie, CSRF, CAPTCHA, redirect, retry, polling, and error handling
   only when the evidence supports the change.
5. Sanitize credentials and account-specific values before committing
   documentation.
6. State what remains unknown instead of presenting a single incidental
   response as a stable public contract.

Configuration data returned by the upstream service, such as language IDs and
record-status IDs, should be resolved at runtime. Do not treat the values in an
example as permanent enumerations.

## Commit Messages

All commits must follow
[Conventional Commits 1.0.0](https://www.conventionalcommits.org/en/v1.0.0/).

```text
<type>[optional scope][optional !]: <description>

[optional body]

[optional footer(s)]
```

The description is required, must be written in imperative English, must not
end with a period, and should describe one logical change.

### Supported types

| Type | Purpose |
|---|---|
| `feat` | Add a user-visible capability |
| `fix` | Correct a defect |
| `docs` | Change documentation only |
| `style` | Change formatting without changing behavior |
| `refactor` | Restructure code without adding a feature or fixing a defect |
| `perf` | Improve performance |
| `test` | Add or correct tests |
| `build` | Change the build system or dependencies |
| `ci` | Change CI configuration |
| `chore` | Perform maintenance that does not fit another type |
| `revert` | Revert an earlier commit |

Use a concise scope when it identifies the affected area:

```text
docs(api): clarify CAPTCHA retry behavior
fix(auth): preserve the final C3VK value
test(record): cover pending judge statuses
```

### Breaking changes

Add `!` before the colon and include a `BREAKING CHANGE:` footer:

```text
refactor(client)!: require an explicit session store

BREAKING CHANGE: callers must pass a session store instead of using the
process-wide cookie jar.
```

The footer must describe the old behavior, the new behavior, and the migration
action.

## Review Checklist

- [ ] Documentation and source comments are in English.
- [ ] API claims are supported by evidence or explicitly labeled unverified.
- [ ] Examples contain no credentials or account-specific values.
- [ ] Runtime configuration is used for language and record-status IDs.
- [ ] Relative links resolve and code blocks are valid.
- [ ] The change updates documentation with the behavior it changes.
- [ ] The commit message follows Conventional Commits.
- [ ] Breaking changes include both `!` and a `BREAKING CHANGE:` footer.
