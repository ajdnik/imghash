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

## Automated Security Scanning

This project employs several automated tools to catch vulnerabilities early:

- **[govulncheck](https://pkg.go.dev/golang.org/x/vuln/cmd/govulncheck)** runs
  in CI on every pull request and push to `main` to catch known vulnerabilities
  in dependencies.
- **[CodeQL](https://codeql.github.com/)** performs static application security
  testing (SAST) on every pull request, push to `main`, and on a weekly schedule.
- **[Dependabot](https://docs.github.com/en/code-security/dependabot)** monitors
  Go module dependencies and GitHub Actions for available updates and opens pull
  requests automatically.
