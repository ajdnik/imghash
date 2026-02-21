# Security Policy

## Supported Versions

Only the latest minor release receives security updates.

| Version | Supported |
|---------|-----------|
| latest | Yes |
| older | No |

## Reporting a Vulnerability

**Please do not open a public issue for security vulnerabilities.**

Instead, report them privately through
[GitHub Security Advisories](https://github.com/ajdnik/imghash/security/advisories/new).

Include as much of the following as you can:

- A description of the vulnerability and its potential impact
- Steps to reproduce or a proof of concept
- Affected versions
- Any suggested fix, if you have one

You should receive an initial acknowledgement within 72 hours. From there the
maintainers will work with you to understand the issue, confirm it, and
coordinate a fix and disclosure timeline.

## Dependency Scanning

This project runs [govulncheck](https://pkg.go.dev/golang.org/x/vuln/cmd/govulncheck)
in CI on every pull request and push to `main` to catch known vulnerabilities in
dependencies.
