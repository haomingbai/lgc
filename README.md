# LGC

LGC is a Luogu code submission project designed to let an LLM submit code to
Luogu and inspect the judging result. It targets workflows where code boards or
solution snippets need to be tested against real problems rather than only
local unit tests.

The intended workflow is:

1. Authenticate with the Luogu web service while keeping the human CAPTCHA step
   under user control.
2. Open a problem and resolve its runtime submission state.
3. Resolve the correct language ID from Luogu's runtime configuration.
4. Submit source code and receive a judge record ID.
5. Wait for asynchronous judging through record polling or WebSocket updates.
6. Return a machine-readable verdict, diagnostics, and test progress to the
   caller.

## Project status

This repository currently contains API research and documentation. It does not
yet include a CLI or library implementation, so no installation or invocation
command is documented here.

All project documentation and contribution requirements are maintained under
[`docs/`](docs/).
