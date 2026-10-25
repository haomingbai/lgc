# Luogu Upstream API Reference

This document describes the Luogu web API used by the project to authenticate,
open a problem, submit source code, and retrieve judge results. The base URL is
`https://www.luogu.com.cn`.

The API is a browser-facing service rather than a stable public API. Its
headers, cookies, routes, language IDs, and challenge behavior can change
without notice. Resolve mutable values at runtime whenever possible.

## Evidence status

Claims in this reference use three status labels:

| Status | Meaning |
|---|---|
| **Verified** | Directly observed in the captured request/response or runtime configuration described in the evidence section. |
| **Inferred** | Consistent with the capture or client behavior but not exercised end-to-end in the available evidence. |
| **Unverified** | Reproduced only from an analysis note or historical client behavior; revalidate before relying on it. |

The capture demonstrates one successful password-login flow, one source-code
submission, one pending record read, and one WebSocket record update. It does
not contain a successful `/api/verify/captcha` request, an answer-file
submission, a contest submission, or an interactive CAPTCHA exchange.

## Verified request flow

```text
GET  /auth/login
GET  /_lfe/config/auth?_version=<config-version>
GET  /auth/login-methods?login=<account>
GET  /lg4/captcha?_t=<timestamp>              -> 302 self-redirect
GET  /lg4/captcha?_t=<same-timestamp>         -> image/jpeg
POST /do-auth/password                        -> session cookies
GET  /problem/<pid>                           -> 302 self-redirect
GET  /problem/<pid>                           -> final HTML and CSRF token
POST /fe/api/problem/submit/<pid>             -> {"rid": <number>}
GET  /record/<rid>                            -> current record JSON
GET  wss://ws2.luogu.com.cn/ws                -> record update
```

The first redirect in each challenge flow is part of the contract observed in
the capture. Follow redirects with a cookie jar enabled and retain cookies set
by the final response.

## Endpoint index

| Method | Path | Evidence | Purpose |
|---|---|---|---|
| `GET` | `/auth/login` | Verified | Load the login page and obtain a CSRF token |
| `GET` | `/auth/login-methods` | Verified | Discover available login methods |
| `GET` | `/lg4/captcha` | Verified | Obtain the login image CAPTCHA and rotate `C3VK` |
| `POST` | `/do-auth/password` | Verified | Authenticate with a password and CAPTCHA |
| `GET` | `/problem/{pid}` | Verified | Load a problem and obtain page-scoped state |
| `POST` | `/fe/api/problem/submit/{pid}` | Verified | Submit source code for asynchronous judging |
| `GET` | `/record/{rid}` | Verified | Read the current judge record |
| `GET` | `/_lfe/config/auth` | Verified | Resolve routes and runtime numeric IDs |
| `GET` | `wss://ws2.luogu.com.cn/ws` | Verified | Receive record updates over WebSocket |
| `GET` | `/api/verify/captcha` | Unverified | Retrieve a submission CAPTCHA after risk control |

## Shared state and headers

### Cookies

Maintain one cookie jar for `www.luogu.com.cn` and send its applicable cookies
on main-site requests. Important cookie names observed in the capture are:

| Name | Observed behavior |
|---|---|
| `_uid` | Anonymous before login and authenticated after login |
| `__client_id` | Session cookie; rotates after successful password login |
| `C3VK` | Short-lived page/challenge cookie; `Max-Age=300` in `Set-Cookie` responses |

The successful login response sets both `_uid` and `__client_id` with
`Secure`, `HttpOnly`, `SameSite=None`, `Domain=.luogu.com.cn`, and
`Max-Age=2592000`. Treat all cookie values as opaque. Do not serialize them
into logs, source code, documentation, or bug reports.

The capture shows `C3VK` changing during each two-request challenge flow.
After following a redirect, use the latest value from the final response.

### CSRF

The CSRF token is not a cookie. Read it from the login or problem HTML:

```html
<meta name="csrf-token" content="<csrf-token>">
```

Send it as:

```http
X-CSRF-Token: <csrf-token>
```

The capture uses the login-page token for login-method discovery and password
login. It later uses the final problem-page token for submission and record
reading. Refresh the token after reloading the corresponding HTML page.

Do not confuse `X-CSRF-Token` with `syncToken`, `C3VK`, or `__client_id`.

### Common request headers

The following headers were observed on JSON requests:

```http
Accept: application/json, text/plain, */*
X-Requested-With: XMLHttpRequest
Origin: https://www.luogu.com.cn
```

Use the relevant page URL as `Referer`. These headers are useful for
compatibility with the web client, but the capture does not prove that every
header is independently required by the server.

## Runtime configuration

## `GET /_lfe/config/auth`

**Evidence:** Verified in frames 3323 (request), 3341 (response headers), and
3343 (decompressed JSON).

### Request

```http
GET /_lfe/config/auth?_version=<config-version> HTTP/2
Host: www.luogu.com.cn
Accept: application/json, text/plain, */*
X-Requested-With: XMLHttpRequest
```

The captured page supplied the `_version` value. Whether the endpoint works
without that query parameter is unverified.

### Response

The response is JSON compressed with gzip. Relevant top-level keys are:

| Key | Meaning |
|---|---|
| `route` | Named endpoint path templates |
| `CodeLanguage` | Submission-language configuration keyed by numeric ID |
| `RecordStatus` | Judge-status configuration keyed by numeric ID |

Selected verified route values:

```json
{
  "auth.login": "/auth/login",
  "auth.login_methods": "/auth/login-methods",
  "do_auth.password": "/do-auth/password",
  "captcha": "/lg4/captcha",
  "problem.show": "/problem/{pid}",
  "api.problem.submit": "/fe/api/problem/submit/{pid}",
  "api.verify.captcha": "/api/verify/captcha",
  "record.show": "/record/{id}"
}
```

Use the numeric `CodeLanguage.*.id` field as the submission `lang`. Reject
languages whose `disabled` field is `true`. Do not hard-code the mapping.

Selected language values observed in this configuration:

| ID | Name | O2 supported |
|---:|---|:---:|
| `27` | C++20 | Yes |
| `34` | C++23 | Yes |
| `7` | Python 3 | No |
| `15` | Rust | Yes |
| `14` | Go | No |
| `8` | Java 8 | No |
| `33` | Java 21 | No |
| `9` | Node.js LTS | No |
| `5` | Answer file | No |

The complete configuration contained 35 entries, 27 enabled and 8 disabled.
These counts are observations from the capture, not permanent guarantees.

## Authentication

## `GET /auth/login`

**Evidence:** Verified in frames 3135, 3151, and 3152.

### Request

```http
GET /auth/login HTTP/2
Host: www.luogu.com.cn
Accept: text/html
Cookie: <current-cookie-jar>
```

### Response

The captured response was HTTP `200`, `text/html; charset=UTF-8`, gzip
compressed, and contained the CSRF meta element described above.

An HTML response without the CSRF element cannot be used reliably for the
captured login flow. Re-fetch the page instead of guessing a token.

## `GET /auth/login-methods`

**Evidence:** Verified in frames 3822 (request), 3837 (response headers), and
3839 (decompressed JSON).

### Request

```http
GET /auth/login-methods?login=<url-encoded-account> HTTP/2
Host: www.luogu.com.cn
Accept: application/json, text/plain, */*
X-Requested-With: XMLHttpRequest
X-CSRF-Token: <login-csrf-token>
Referer: https://www.luogu.com.cn/auth/login
Cookie: <current-cookie-jar>
```

URL-encode the account identifier. The captured identifier was a username;
the supported input forms for UID, phone number, and email are **Unverified**.

### Response

The observed account supported password authentication:

```json
{
  "available": ["password"],
  "default": "password"
}
```

Only the password path is verified end-to-end. Route configuration mentions
FIDO2, TOTP, and MOTP endpoints, but their request and response contracts are
**Unverified**.

## `GET /lg4/captcha`

**Evidence:** Verified in frames 3841 and 3843 (first request and `302`), then
3847 and 3855 (replay and final image).

### Request

```http
GET /lg4/captcha?_t=<millisecond-timestamp> HTTP/2
Host: www.luogu.com.cn
Accept: image/avif,image/jxl,image/webp,image/png,image/svg+xml,image/*;q=0.8,*/*;q=0.5
Referer: https://www.luogu.com.cn/auth/login
Cookie: <current-cookie-jar>
```

The captured `_t` value contained a fractional millisecond timestamp. A
simple integer millisecond timestamp is reasonable, but the no-fraction form
is **Inferred**.

### Redirect response

The first response was:

```http
HTTP/2 302
Location: https://www.luogu.com.cn/lg4/captcha?_t=<same-timestamp>
Content-Type: text/html
Set-Cookie: C3VK=<first-value>; Max-Age=300; Path=/
```

Follow the redirect with cookies enabled. The captured client replayed the same
URL after storing the first `C3VK` value.

### Final response

```http
HTTP/2 200
Content-Type: image/jpeg
Set-Cookie: C3VK=<latest-value>; Max-Age=300; Path=/
```

Present the final image to a human and use the latest `C3VK` cookie in the
password request. CAPTCHA answers are single-use.

## `POST /do-auth/password`

**Evidence:** Verified in frames 4177 (request), 4190 (response headers), and
4192 (decompressed JSON).

### Request

```http
POST /do-auth/password HTTP/2
Host: www.luogu.com.cn
Content-Type: application/json
Accept: application/json, text/plain, */*
X-Requested-With: XMLHttpRequest
X-CSRF-Token: <login-csrf-token>
Origin: https://www.luogu.com.cn
Referer: https://www.luogu.com.cn/auth/login
Cookie: C3VK=<latest-login-c3vk>; _uid=<anonymous-uid>; __client_id=<anonymous-client-id>
```

```json
{
  "username": "<account>",
  "password": "<password>",
  "captcha": "<human-entered-answer>"
}
```

| Field | Type | Meaning |
|---|---|---|
| `username` | string | Account identifier |
| `password` | string | Account password |
| `captcha` | string | Text entered for the current image challenge |

Do not log or persist the password or CAPTCHA answer.

### Success response

The captured response was HTTP `200` with gzip-compressed JSON:

```json
{
  "username": "<account>",
  "locked": false,
  "syncToken": "<opaque-token>",
  "redirectTo": "https://www.luogu.com.cn/"
}
```

The response also set rotated `_uid` and `__client_id` cookies. Treat login as
successful only when the expected JSON is present and both session cookies are
stored. The captured `syncToken` was not sent by subsequent problem or
submission requests.

### Failure behavior

The error examples and validation order in the pre-review document were not
present in the capture. Specific exception classes, messages, and HTTP status
codes are therefore **Unverified**.

On any malformed or error response, discard the current CSRF and CAPTCHA state,
then restart the login flow if another manual attempt is appropriate. Do not
retry a consumed CAPTCHA automatically.

## Problem and submission

## `GET /problem/{pid}`

**Evidence:** Verified in frames 4956 and 4958 (first request and `302`), then
4960, 5004, and the following response data (replay and final HTML).

### Request

```http
GET /problem/<pid> HTTP/2
Host: www.luogu.com.cn
Accept: text/html
Referer: https://www.luogu.com.cn/
Cookie: <authenticated-cookie-jar>
```

`pid` is an opaque problem identifier such as `P1001`.

### Redirect and final response

```http
HTTP/2 302
Location: https://www.luogu.com.cn/problem/<pid>
Set-Cookie: C3VK=<first-value>; Max-Age=300; Path=/
```

After following the redirect with the cookie jar:

```http
HTTP/2 200
Content-Type: text/html; charset=UTF-8
Set-Cookie: C3VK=<latest-value>; Max-Age=300; Path=/
```

The final HTML contained a CSRF meta token. Use the latest `C3VK` and final
page token for the subsequent submission and record request.

## `POST /fe/api/problem/submit/{pid}`

**Evidence:** Verified in frames 7387 (request), 7426 (response headers), and
7428 (decompressed JSON).

### Source-code request

```http
POST /fe/api/problem/submit/<pid> HTTP/2
Host: www.luogu.com.cn
Content-Type: application/json
Accept: application/json, text/plain, */*
X-Requested-With: XMLHttpRequest
X-CSRF-Token: <problem-csrf-token>
Origin: https://www.luogu.com.cn
Referer: https://www.luogu.com.cn/problem/<pid>
Cookie: C3VK=<latest-problem-c3vk>; _uid=<authenticated-uid>; __client_id=<authenticated-client-id>
```

```json
{
  "lang": 27,
  "code": "#include <iostream>\nint main() { return 0; }\n",
  "enableO2": 1
}
```

| Field | Type | Meaning |
|---|---|---|
| `lang` | integer | Runtime `CodeLanguage.*.id` for the selected language |
| `code` | string | Complete source text |
| `enableO2` | integer | `1` to request O2; `0` otherwise |

The capture used `lang: 34` for C++23. The example above uses `27` only
because C++20 is enabled in the captured configuration.

### Success response

```http
HTTP/2 200
Content-Type: application/json
Content-Encoding: gzip
```

```json
{
  "rid": 123456789
}
```

`rid` is a numeric record identifier in the captured response. Treat it as an
opaque identifier when constructing `/record/<rid>`. The submission response
does not contain the final verdict; judging is asynchronous.

### Unverified request variants

The pre-review document described these variants, but the capture does not
exercise them:

| Variant | Status | Caveat |
|---|---|---|
| `multipart/form-data` answer-file upload | **Unverified** | Field names and accepted file types were not observed |
| `?contestId=<id>` query parameter | **Unverified** | Contest routing and required context were not exercised |
| JSON `captcha` retry field | **Unverified** | See the challenge section below |
| Interactive CAPTCHA token format | **Unverified** | Do not implement from the historical description without a new capture |

Do not claim support for these variants until a sanitized request and response
have been captured and added as evidence.

## Judge records

## `GET /record/{rid}`

**Evidence:** Verified in frames 7537 (request), 7558 (response headers), and
7559 (decompressed JSON).

### Request

```http
GET /record/<rid> HTTP/2
Host: www.luogu.com.cn
Accept: application/json, text/plain, */*
X-Requested-With: XMLHttpRequest
X-Lentille-Request: content-only
Referer: https://www.luogu.com.cn/problem/<pid>
X-CSRF-Token: <problem-csrf-token>
Cookie: C3VK=<problem-c3vk>; _uid=<authenticated-uid>; __client_id=<authenticated-client-id>
```

The captured record belonged to the authenticated account and was requested
from the corresponding problem page. The capture does not establish whether
public records can be read with a different authentication or header set.

### Response

The response was HTTP `200` with gzip-compressed JSON:

```json
{
  "instance": "main",
  "template": "record.show",
  "status": 200,
  "data": {
    "record": {
      "id": 123456789,
      "status": 0,
      "language": 34,
      "sourceCodeLength": 154,
      "detail": {
        "compileResult": {
          "success": false,
          "message": "<compiler-output>",
          "opt2": false
        },
        "judgeResult": {
          "finishedCaseCount": 0,
          "status": 0,
          "time": 0,
          "memory": 0,
          "score": 0,
          "subtasks": []
        }
      }
    }
  }
}
```

Important distinctions:

- Top-level `status` is the Lentille envelope status.
- `data.record.status` is the current judge status.
- `data.record.detail.compileResult` contains compilation state and diagnostics.
- `data.record.detail.judgeResult` contains case progress, score, and subtasks.
- `data.record.problem.accepted` describes the user/problem relationship and
  is not the verdict for this record.

While `data.record.status` is `0`, the record is still waiting. The capture
shows `compileResult.success: false` even for this pending, non-terminal
record. Do not classify the record as Compile Error until the overall record
status is terminal.

Poll with a bounded exponential backoff and stop once `data.record.status` is
not `0` or `1`. The capture only directly observes `0` over HTTP and `2` over
WebSocket; the rest of the status mapping comes from runtime configuration.

## Record status configuration

**Evidence:** Verified in the decompressed config response, frame 3343.

| ID | Type | Short name | Poll again |
|---:|---|---|:---:|
| `-1` | Unshown | N/A | No |
| `0` | Waiting | Waiting | Yes |
| `1` | Judging | Judging | Yes |
| `2` | Compile Error | CE | No |
| `3` | Output Limit Exceeded | OLE | No |
| `4` | Memory Limit Exceeded | MLE | No |
| `5` | Time Limit Exceeded | TLE | No |
| `6` | Wrong Answer | WA | No |
| `7` | Runtime Error | RE | No |
| `11` | Invalid / Unknown Error | UKE | No |
| `12` | Accepted | AC | No |
| `14` | Unaccepted | NAC | No |
| `21` | Hack Success | Hack Success | No |
| `22` | Hack Failure | Hack Failure | No |
| `23` | Hack Skipped | Hack Skipped | No |

The `Poll again` column is **Inferred** from the configuration and the
single-record behavior. Resolve status names from a fresh
`GET /_lfe/config/auth` response at runtime.

## WebSocket updates

## `wss://ws2.luogu.com.cn/ws`

**Evidence:** Verified in frames 5265 and 5295 (HTTP/1.1 upgrade to `101`) and
frames 7591 and 7602 (record update).

### Handshake

The captured request uses HTTP/1.1 WebSocket upgrade:

```http
GET /ws HTTP/1.1
Host: ws2.luogu.com.cn
Upgrade: websocket
Connection: Upgrade
Origin: https://fecdn.luogu.com.cn
Sec-WebSocket-Version: 13
Sec-WebSocket-Extensions: permessage-deflate
Cookie: <authenticated-cookie-jar>
```

The server returned `101 Switching Protocols`. The exact set of cookies
needed for a standalone client is not established by one capture; reuse the
authenticated jar supported for `luogu.com.cn`.

### Record update frame

After decompression, the observed binary payload has this structure:

```text
0x10 + "Record:<rid>" + 0x01 + <record-json>
```

The JSON is a record object using the same judge-status fields described
above. The captured update changed the record status from waiting (`0`) to
Compile Error (`2`).

This capture does not define subscription messages, reconnect semantics,
heartbeats, or delivery guarantees. Use `GET /record/{rid}` as the canonical
snapshot and polling fallback. On reconnect, poll the record before relying on
future pushes.

## Submission CAPTCHA branch

**Status: Unverified.** The endpoint name exists in the runtime route map, but
the capture contains no request to `/api/verify/captcha` and no submission
response that requests a CAPTCHA.

The historical analysis describes this possible path:

1. A submission can fail with CAPTCHA-related JSON.
2. The client fetches an image from `/api/verify/captcha`.
3. The endpoint can issue a self-redirect while setting `C3VK`.
4. The client resubmits the original JSON body with a `captcha` field.

This behavior is plausible but not proven by the available artifacts. In
particular, the following details are **Unverified**:

- exact HTTP status and error envelope;
- whether the timestamp query parameter is required;
- retry-count limits;
- interactive challenge providers and token encoding;
- signature/data fields used for an interactive provider;
- whether a `403` response is returned before or after challenge state is set.

The client should detect the branch only from a live response, avoid unbounded
retries, and require human interaction for an interactive challenge. Until a
sanitized capture is added, do not advertise answer-file, contest, or
interactive-CAPTCHA submission as supported features.

## Error handling

HTTP success alone is not enough to prove an API operation succeeded:

- Parse JSON when the endpoint normally returns JSON.
- Validate endpoint-specific fields such as a non-zero `rid` or a record object.
- Check root-level `errorCode`, `errorType`, and `errorMessage` if present.
- Check Lentille responses for `template: "error"` and an error envelope status.
- Stop on malformed or unexpected content and surface a diagnostic.

The generic envelopes in the pre-review document are **Unverified** because the
capture contains no failing JSON response for login, submission, or record
retrieval. Do not match on user-facing error messages. If matching exceptions,
use exact documented identifiers only after adding evidence for them.

## Evidence

The primary evidence is:

- `.artifacts/packet1.pcapng`
- `.artifacts/ssl-key.log`
- `.artifacts/packet1.md`

Capture metadata:

| Property | Value |
|---|---|
| Packets | 8956 |
| Capture duration | 105.051 seconds |
| First packet | 2026-10-24 21:43:15 (UTC+8) |
| Last packet | 2026-10-24 21:45:00 (UTC+8) |
| Protocol | TLS with HTTP/2 and HTTP/1.1 WebSocket |

The capture and TLS key log can decrypt account-specific traffic. They are
local evidence artifacts and must not be published. The repository ignores
`.artifacts/`.
