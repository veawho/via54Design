# Security Policy

## Supported Versions

| Version | Supported          |
| ------- | ------------------ |
| v0.4.x  | ✅ active          |
| v0.3.x  | ✅ security fixes  |
| v0.2.x  | ⚠️ best-effort     |
| < v0.2  | ❌ end-of-life     |

## Reporting a Vulnerability

If you discover a security vulnerability in `via54Design`, **please do not
open a public GitHub issue**. Instead, use one of these private channels:

1. **GitHub Security Advisories** (preferred):  
   https://github.com/veawho/via54Design/security/advisories/new

2. **Email**: `veawho@protonmail.com`  
   Subject line: `[via54Design security]`

Please include the following information in your report:

- A clear description of the vulnerability
- Steps to reproduce (or a proof-of-concept)
- The affected version(s) and commit hash
- The potential impact (e.g., arbitrary code execution, information disclosure)
- Your assessment of severity (Critical / High / Medium / Low)

We will acknowledge receipt within **72 hours** and aim to release a fix
within **30 days** for high-severity issues.

## Scope

The following components are in scope:

- Go binaries (`via54`, `via54-mcp`) — input validation, path traversal, command injection
- Web UI (HTMX endpoints) — file upload, RCE via template injection
- MCP Server — tool argument injection, path traversal
- YAML template loader — type confusion, prototype-style attacks
- Export pipeline — ZIP/XML parsing, external command invocation
- All third-party Go dependencies listed in `go.mod`

## Out of scope

- Issues in the user's local environment (broken Node.js, misconfigured ffmpeg)
- Vulnerabilities in **upstream** projects (HTMX, mcp-go, yaml.v3) — please report those upstream
- Denial-of-service from extremely large `--scene` / `--seed` strings
  (we have length limits but they are tunable, not security boundaries)
- Social engineering attacks against the maintainer

## Coordinated disclosure

We follow a **90-day disclosure timeline**:

- Day 0: vulnerability reported
- Day 0–7: triage and confirm
- Day 7–60: develop and test a fix
- Day 60–90: coordinate disclosure with reporter
- Day 90: public advisory + CVE assigned (if applicable)

If a fix requires more time, we will negotiate an extension with the
reporter before the 90-day deadline.

## Recognition

We maintain a list of security researchers who have helped us in
[`ACKNOWLEDGMENTS.md`](./ACKNOWLEDGMENTS). With your permission, we will
add your name there.

Thank you for keeping via54Design — and its users — safe.
