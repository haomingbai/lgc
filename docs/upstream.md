# Luogu Upstream API Reference

This document describes the Luogu web API used by the project to authenticate,
optionally read problem content, submit source code, and retrieve judge
results. The base URL is `https://www.luogu.com.cn`.

The API is a browser-facing service rather than a stable public API. Its
headers, cookies, routes, language IDs, and challenge behavior can change
without notice. Resolve mutable values at runtime whenever possible.

## Evidence status

Claims in this reference use three status labels:

| Status | Meaning |
|---|---|
| **Verified** | Directly observed in the captured request/response, runtime configuration, or controlled live test described in the evidence section. |
| **Inferred** | Consistent with the capture or client behavior but not exercised end-to-end in the available evidence. |
| **Unverified** | Reproduced only from an analysis note or historical client behavior; revalidate before relying on it. |

The original capture demonstrates one successful password-login flow, one
source-code submission, one pending record read, and one WebSocket record
update. A later live test exercised the submission CAPTCHA branch and showed
that a submission does not require a problem-page request.

Neither evidence source demonstrates an answer-file submission, a contest
submission, or a provider-specific interactive challenge exchange.

## Observed request flows

### Captured browser flow

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

This sequence records the browser's UI behavior. The problem-page request is
not a submission dependency.

### Live-tested direct-submission flow

```text
GET  /                                        -> authenticated CSRF token
POST /fe/api/problem/submit/<pid>             -> 403 CAPTCHA required
GET  /api/verify/captcha                      -> image challenge
POST /fe/api/problem/submit/<pid>             -> {"rid": <number>}
GET  /record/<rid>                            -> current record JSON
```

The live test never requested `/problem/<pid>` or another problem page. The
retry POST was the original submission body plus one `captcha` field.

This test triggered risk-control CAPTCHA. That challenge is conditional: the
original browser capture submitted without one. Clients should attempt the base
POST first and enter the CAPTCHA branch only when the response requests it.

A `302` self-redirect may occur while obtaining a page or challenge. Follow it
with a cookie jar enabled only when it is present; a redirect is not a universal
precondition.

## Endpoint index

| Method | Path | Evidence | Data provided by the endpoint |
|---|---|---|---|
| `GET` | `/` | Verified | Authenticated HTML and a CSRF meta token |
| `GET` | `/auth/login` | Verified | Login HTML and a login CSRF meta token |
| `GET` | `/auth/login-methods` | Verified | Available login method names and the default method |
| `GET` | `/lg4/captcha` | Verified | Login CAPTCHA image and rotating challenge cookies |
| `POST` | `/do-auth/password` | Verified | Authenticated session cookies and login JSON |
| `GET` | `/_lfe/config/auth` | Verified | Route templates, language IDs, and record-status IDs |
| `GET` | `/problem/{pid}` | Verified | Problem HTML and optional problem content |
| `POST` | `/fe/api/problem/submit/{pid}` | Verified | A RID, or a CAPTCHA challenge error |
| `GET` | `/api/verify/captcha` | Verified | Submission CAPTCHA image and any updated cookies |
| `GET` | `/record/{rid}` | Verified | Current judge-record snapshot |
| `GET` | `wss://ws2.luogu.com.cn/ws` | Verified | Asynchronous record updates |

### Minimal data-source map

This table identifies the smallest observed source for each value used by the
submission workflow. "Caller" means data supplied to the client rather than
obtained from Luogu.

| Requirement | Source | When needed |
|---|---|---|
| Problem ID | Caller | Every submission |
| Authenticated cookies | `POST /do-auth/password` response | Every authenticated request |
| CSRF token | `GET /` HTML in the live test | Authenticated JSON requests |
| Language ID | `GET /_lfe/config/auth` response | Resolve the requested language |
| Source text and O2 choice | Caller | Every submission |
| First submission result | `POST /fe/api/problem/submit/{pid}` response | Always |
| Submission CAPTCHA | `GET /api/verify/captcha`, only after a CAPTCHA error | When the POST returns the verified challenge error |
| RID | Successful submission response | Waiting and record retrieval |
| Record snapshot | `GET /record/{rid}` | Until a terminal status |
| Low-latency record updates | WebSocket | Optional; must not replace polling fallback |
| Problem content | `GET /problem/{pid}` | Only when the caller requests content |

### Endpoint parameter sources

This table is the authoritative cross-endpoint reference for data-bearing
request parameters. Protocol constants such as `Host`, `Origin`, and static
`Accept` values do not require another endpoint.

| Endpoint | Parameter | Source or reference |
|---|---|---|
| `GET /` | `Cookie` | Current cookie jar; authenticated `_uid` and `__client_id` come from `POST /do-auth/password` |
| `GET /auth/login` | `Cookie` | Current cookie jar, including cookies established by preceding site navigation |
| `GET /_lfe/config/auth` | `_version` query | The HTML/client runtime loaded from a page such as `GET /auth/login`; this endpoint does not return its own `_version` |
| `GET /_lfe/config/auth` | `Cookie` | Current cookie jar |
| `GET /auth/login-methods` | `login` query | Caller-provided account identifier |
| `GET /auth/login-methods` | `X-CSRF-Token` | The CSRF meta token from `GET /auth/login` |
| `GET /auth/login-methods` | `Cookie` | Current cookie jar or the jar stored after `GET /auth/login` |
| `GET /auth/login-methods` | `Referer` | `https://www.luogu.com.cn/auth/login` |
| `GET /lg4/captcha` | `_t` query | Client-generated timestamp |
| `GET /lg4/captcha` | `Cookie` | Current cookie jar; preserve cookies set by each challenge response |
| `GET /lg4/captcha` | `Referer` | `https://www.luogu.com.cn/auth/login` |
| `POST /do-auth/password` | `username`, `password` | Caller-provided credentials |
| `POST /do-auth/password` | `captcha` | The image returned by `GET /lg4/captcha`, processed by the CAPTCHA solver |
| `POST /do-auth/password` | `X-CSRF-Token` | The CSRF meta token from `GET /auth/login` |
| `POST /do-auth/password` | `C3VK` | The latest `C3VK` cookie from the `GET /lg4/captcha` challenge flow |
| `POST /do-auth/password` | `_uid`, `__client_id` | The current anonymous cookie jar populated by preceding main-site responses |
| `POST /do-auth/password` | `Referer` | `https://www.luogu.com.cn/auth/login` |
| `GET /problem/{pid}` | `pid` | Caller or an existing `Problem` object |
| `GET /problem/{pid}` | `Cookie` | Authenticated cookie jar from `POST /do-auth/password` |
| `GET /problem/{pid}` | `Referer` | `https://www.luogu.com.cn/` |
| `POST /fe/api/problem/submit/{pid}` | `pid` | Caller or the `Problem` object being submitted |
| `POST /fe/api/problem/submit/{pid}` | `X-CSRF-Token` | Authenticated CSRF meta token from `GET /` |
| `POST /fe/api/problem/submit/{pid}` | `Cookie` | Authenticated cookie jar from `POST /do-auth/password` |
| `POST /fe/api/problem/submit/{pid}` | `Referer` | `https://www.luogu.com.cn/` in the live test |
| `POST /fe/api/problem/submit/{pid}` | `lang` | `CodeLanguage.*.id` from `GET /_lfe/config/auth` |
| `POST /fe/api/problem/submit/{pid}` | `code` | Caller-provided source text |
| `POST /fe/api/problem/submit/{pid}` | `enableO2` | Caller-provided O2 choice |
| `POST /fe/api/problem/submit/{pid}` | `captcha` | Only after a CAPTCHA error: the image from `GET /api/verify/captcha`, processed by the CAPTCHA solver |
| `GET /api/verify/captcha` | Trigger | The preceding submission response from `POST /fe/api/problem/submit/{pid}` |
| `GET /api/verify/captcha` | `_t` query | Client-generated timestamp |
| `GET /api/verify/captcha` | `X-CSRF-Token` | Authenticated CSRF meta token from `GET /` |
| `GET /api/verify/captcha` | `Cookie` | Authenticated cookie jar plus cookies updated by the submission CAPTCHA flow |
| `GET /api/verify/captcha` | `Referer` | `https://www.luogu.com.cn/` in the live test |
| `GET /record/{rid}` | `rid` | The successful `POST /fe/api/problem/submit/{pid}` response |
| `GET /record/{rid}` | `X-CSRF-Token` | Authenticated CSRF meta token from `GET /` |
| `GET /record/{rid}` | `Cookie` | Authenticated cookie jar from `POST /do-auth/password` |
| `GET /record/{rid}` | `Referer` | `https://www.luogu.com.cn/` in the live test |
| WebSocket | URL | The captured `wss://ws2.luogu.com.cn/ws` endpoint; runtime discovery is unverified |
| WebSocket | `Cookie` | Authenticated cookie jar from `POST /do-auth/password` |
| WebSocket | `Origin` | Static web-client value `https://fecdn.luogu.com.cn` |

The exclusive source of an authenticated CSRF token has not been established.
The live test proves that the root page works and that the problem page is not
required.

Response fields are provided by the endpoint documented in the section where
they appear; they are not obtained from a different Luogu API unless the text
explicitly says so.

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

`C3VK` can rotate when a challenge endpoint issues `Set-Cookie`. Preserve the
latest value in the session cookie jar. The live submission CAPTCHA request
returned the image directly, so a two-request redirect flow is possible but not
required for every challenge.

### CSRF

The CSRF token is not a cookie. Read it from an HTML `<meta>` element:

```html
<meta name="csrf-token" content="<csrf-token>">
```

Send it as:

```http
X-CSRF-Token: <csrf-token>
```

The capture uses the login-page token for login-method discovery and password
login. The live test read a fresh token from `GET /` and used it for submission,
the submission CAPTCHA request, and record polling without loading the problem
page.

Refresh the token if it is missing or rejected. Do not fetch a problem page
solely to obtain one.

Do not confuse `X-CSRF-Token` with `syncToken`, `C3VK`, or `__client_id`.

### Common request headers

The following headers were observed on JSON requests:

```http
Accept: application/json, text/plain, */*
X-Requested-With: XMLHttpRequest
Origin: https://www.luogu.com.cn
```

The captured browser used the current problem page as `Referer`. The live test
used `https://www.luogu.com.cn/` for direct submission, submission CAPTCHA, and
record polling. These headers matched the tested client behavior, but no
evidence establishes that every header is independently required.

## `GET /`

**Evidence:** Verified by the controlled live submission test on 2026-10-26.

### Request

```http
GET / HTTP/2
Host: www.luogu.com.cn
Accept: text/html
Cookie: <authenticated-cookie-jar>
```

### Response

The observed response was HTTP `200` with HTML containing the CSRF `<meta>`
element. In the live test, the token was used for direct submission, the
submission CAPTCHA request, and record polling.

This endpoint provides an authenticated HTML page and CSRF token. It does not
provide problem content. The test does not establish that the root page is the
only source of a usable authenticated CSRF token.

## Runtime configuration

## `GET /_lfe/config/auth`

**Evidence:** Verified by the runtime configuration response observed in the
browser capture on 2026-10-24.

### Request

```http
GET /_lfe/config/auth?_version=<config-version> HTTP/2
Host: www.luogu.com.cn
Accept: application/json, text/plain, */*
X-Requested-With: XMLHttpRequest
```

The captured page's client runtime supplied the `_version` value; it is not
returned by this configuration endpoint. Whether the endpoint works without
that query parameter is unverified.

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

| ID | Name | `canO2` |
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

### O2 capability

Each `CodeLanguage` entry has a boolean `canO2` field. In the captured
configuration, the following entries had `canO2: true`:

| ID | Name |
|---:|---|
| `1` | Pascal |
| `2` | C |
| `3` | C++98 |
| `4` | C++11 |
| `11` | C++14 |
| `12` | C++17 |
| `27` | C++20 |
| `28` | C++14 (GCC 9) |
| `34` | C++23 |
| `15` | Rust |
| `19` | Haskell |
| `30` | OCaml |

Java 8, Java 21, Python 3, PyPy 3, Go, Julia, Lua, Kotlin/JVM, Scala,
C# Mono, Node.js, PHP, Ruby, and Perl had `canO2: false` in the same
configuration. O2 support is therefore not limited to C and C++, nor is it
available for every enabled language.

The client should read `canO2` at runtime and reject an O2 request for a
language whose capability is false. A controlled live submission on 2026-10-26
sent Java 21 with `enableO2: 1`; the submission returned HTTP `200`, but the
terminal record contained `enableO2: false`. This shows that the server may
silently downgrade the request rather than reject it.

The client-side constraint is therefore stronger than the HTTP contract:
unknown, disabled, or `canO2: false` languages are rejected before an O2
submission is sent. The client must not use the server's silent downgrade as a
capability check or treat the resulting record as proof that O2 was enabled.

## Authentication

## `GET /auth/login`

**Evidence:** Verified by the browser capture on 2026-10-24.

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

**Evidence:** Verified by the browser capture on 2026-10-24.

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

The account identifier comes from the caller. The CSRF token and cookies come
from `GET /auth/login` and the current session jar.

## `GET /lg4/captcha`

**Evidence:** Verified by the browser capture on 2026-10-24.

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

This endpoint provides only the login challenge. It does not provide submission
data or depend on a problem page.

## `POST /do-auth/password`

**Evidence:** Verified by the browser capture on 2026-10-24.

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

The username and password come from the caller. The login CSRF token comes from
`GET /auth/login`; the CAPTCHA answer comes from the image returned by
`GET /lg4/captcha`; and the current cookies come from the session jar.

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

## Problem information and submission

## `GET /problem/{pid}`

**Evidence:** Verified by the browser capture on 2026-10-24.

This endpoint provides problem-page data. It is not required by the verified
submission flow and should remain independent from submission.

### Request

```http
GET /problem/<pid> HTTP/2
Host: www.luogu.com.cn
Accept: text/html
Referer: https://www.luogu.com.cn/
Cookie: <authenticated-cookie-jar>
```

`pid` is an opaque problem identifier such as `P1001`.

The caller supplies `pid`. The authenticated cookie jar comes from
`POST /do-auth/password`.

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

The final HTML contained problem content and a CSRF meta token. These values
describe this page response; they are not required inputs to submission. The
live test obtained a usable CSRF token from `GET /` instead and did not request
this endpoint.

The captured P1001 HTML also embedded an `acceptLanguages` array in page data.
Its extraction contract and consistency with other problems are **Inferred**:
the HTML was observed, but no stable public JSON response was exercised. Treat
problem-specific language restrictions as optional cached hints and do not
require a problem-page request for submission.

## `POST /fe/api/problem/submit/{pid}`

**Evidence:** Verified by the browser capture on 2026-10-24, then revalidated by
a controlled live test on 2026-10-26 without any problem-page request.

### Source-code request

```http
POST /fe/api/problem/submit/<pid> HTTP/2
Host: www.luogu.com.cn
Content-Type: application/json
Accept: application/json, text/plain, */*
X-Requested-With: XMLHttpRequest
X-CSRF-Token: <authenticated-csrf-token>
Origin: https://www.luogu.com.cn
Referer: https://www.luogu.com.cn/
Cookie: <authenticated-cookie-jar>
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

Use the runtime `canO2` capability rather than assuming which languages accept
O2. The endpoint can return HTTP `200` for an unsupported O2 request while the
resulting record records `enableO2: false`.

The request combines data from independent sources:

- The caller supplies `pid`, source code, O2 choice, and the requested language.
- `GET /_lfe/config/auth` supplies the runtime language ID for that language.
- `POST /do-auth/password` supplies the authenticated cookies.
- `GET /` supplied the CSRF token in the live test.

No problem-page request is part of this request's verified dependency set.

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

### Remaining unverified request variants

The available evidence does not exercise these variants:

| Variant | Status | Caveat |
|---|---|---|
| `multipart/form-data` answer-file upload | **Unverified** | Field names and accepted file types were not observed |
| `?contestId=<id>` query parameter | **Unverified** | Contest routing and required context were not exercised |
| Interactive CAPTCHA token format | **Unverified** | Do not implement from the historical description without a new capture |

Do not claim support for these variants until a sanitized request and response
have been observed and added as evidence. The JSON `captcha` retry field is
verified separately in the submission CAPTCHA section.

## Judge records

## `GET /record/{rid}`

**Evidence:** Verified by the browser capture on 2026-10-24, then revalidated by
the live submission test on 2026-10-26.

### Request

```http
GET /record/<rid> HTTP/2
Host: www.luogu.com.cn
Accept: application/json, text/plain, */*
X-Requested-With: XMLHttpRequest
X-Lentille-Request: content-only
Referer: https://www.luogu.com.cn/
X-CSRF-Token: <authenticated-csrf-token>
Cookie: <authenticated-cookie-jar>
```

The RID comes from a successful `POST
/fe/api/problem/submit/{pid}` response. The authenticated cookies come from
login, and the CSRF token can come from `GET /`. The live test used those
sources and polled successfully without visiting a problem page.

The capture does not establish whether public records can be read with a
different authentication or header set.

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
- `data.record.enableO2` is the O2 state associated with this record. It can be
  false even when an unsupported language was submitted with `enableO2: 1`.
- `data.record.detail.compileResult` contains compilation state and diagnostics.
- `data.record.detail.judgeResult` contains case progress, score, and subtasks.
- `data.record.problem.accepted` describes the user/problem relationship and
  is not the verdict for this record.

Do not infer language O2 capability from `compileResult.opt2`. That field was
false in observed terminal records for Java 21, Rust, and C++20, despite their
different runtime `canO2` values.

While `data.record.status` is `0`, the record is still waiting. The capture
shows `compileResult.success: false` even for this pending, non-terminal
record. Do not classify the record as Compile Error until the overall record
status is terminal. The live test observed terminal status `12` with
`compileResult.success: true`.

Poll with a bounded exponential backoff and stop once `data.record.status` is
not `0` or `1`. HTTP polling has directly observed `0` and `12`; WebSocket has
directly observed `2`. Resolve the remaining status names from runtime
configuration.

## Record status configuration

**Evidence:** Verified by the runtime configuration response observed in the
browser capture on 2026-10-24.

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

**Evidence:** Verified by the browser capture on 2026-10-24.

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

**Status: Verified.** The controlled live test on 2026-10-26 exercised this
branch for `P1001` without requesting a problem page.

### Trigger from `POST /fe/api/problem/submit/{pid}`

The first submission POST used the normal source-code body and returned HTTP
`403` with JSON whose `errorType` was:

```text
Luogu\Exception\Request\InvalidCaptchaException
```

Do not depend on the localized `data` or `errorMessage` text. The verified
signal is the exact exception identifier together with the HTTP status.

The trigger and CAPTCHA image come from different endpoints:

- `POST /fe/api/problem/submit/{pid}` reports that a challenge is required.
- `GET /api/verify/captcha` provides the image challenge.

The problem page provides neither value.

### Challenge image from `GET /api/verify/captcha`

The live request used the authenticated cookie jar and root-page CSRF token:

```http
GET /api/verify/captcha?_t=<fractional-millisecond-timestamp> HTTP/2
Host: www.luogu.com.cn
Accept: image/avif,image/jxl,image/webp,image/png,image/svg+xml,image/*
X-CSRF-Token: <authenticated-csrf-token>
Referer: https://www.luogu.com.cn/
Cookie: <authenticated-cookie-jar>
```

### Challenge response

The live run returned HTTP `200 image/jpeg` directly, without a redirect.
Preserve any updated cookies from the response.

Earlier client analysis reported a `302` self-redirect; that behavior remains
**Inferred**. If one is present, replay the same challenge URL with the cookie
jar and use the latest cookies. This branch is not required for the direct
`200` response observed in the live test.

### Retry using `POST /fe/api/problem/submit/{pid}`

After solving the image, retry the original submission with the same base body
plus a single `captcha` string:

```json
{
  "lang": 27,
  "code": "<complete-source-text>",
  "enableO2": 1,
  "captcha": "<recognized-image-text>"
}
```

The live retry returned HTTP `200` with a RID and the record subsequently
reached Accepted. CAPTCHA answers are single-use. Limit retries and stop when
the challenge cannot be solved reliably; do not retry an invalid answer in a
loop.

The following details remain **Unverified**:

- whether `_t` is required rather than merely accepted;
- whether a `302` can also occur in other sessions;
- retry-count and rate-limit boundaries;
- interactive challenge providers and token encoding.

The client should enter this branch only after observing the CAPTCHA error and
should not prefetch a submission CAPTCHA for every submission.

## Error handling

HTTP success alone is not enough to prove an API operation succeeded:

- Parse JSON when the endpoint normally returns JSON.
- Validate endpoint-specific fields such as a non-zero `rid` or a record object.
- Check root-level `errorCode`, `errorType`, and `errorMessage` if present.
- Check Lentille responses for `template: "error"` and an error envelope status.
- Stop on malformed or unexpected content and surface a diagnostic.

The HTTP `403` submission error and
`InvalidCaptchaException` identifier are verified by the live test. Other
failure envelopes for login, submission, and record retrieval remain
**Unverified**. Do not match on user-facing error messages. Use only exact
identifiers that have supporting evidence.

## Evidence

Evidence comes from three sanitized observation sets:

- Browser-observed requests and responses from 2026-10-24.
- A controlled live submission test on 2026-10-26.
- Controlled O2 capability probes on 2026-10-26.

Raw traffic, decryption material, credentials, and temporary challenge files
are not part of this documentation.

Browser-observation metadata:

| Property | Value |
|---|---|
| Date | 2026-10-24 (UTC+8) |
| Protocol | TLS with HTTP/2 and HTTP/1.1 WebSocket |

Live-test metadata:

| Property | Value |
|---|---|
| Date | 2026-10-26 (UTC+8) |
| Problem | `P1001` |
| Authentication source | Existing authenticated session |
| Observed statuses | `200` root, `403` submission, `200` CAPTCHA image, `200` submission, `200` record |
| Problem-page requests | `0` |
| Final judge status | `12` (Accepted), 10/10 test cases |
| Protocol | HTTPS |

The additional O2 probes used minimal P1001 submissions in Java 21, Python 3,
Rust, and C++20. The Java 21 record observed with a top-level `enableO2` field
contained `enableO2: false` after an `enableO2: 1` request. Other probe results
do not prove that the compiler applied O2 solely from the record fields
described here. No credentials, CAPTCHA answers, or real RIDs are retained in
this document.

Credentials, cookies, CSRF values, CAPTCHA answers, and real RIDs are omitted
from this document and from all published examples.
