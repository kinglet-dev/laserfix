# Security policy

## Reporting a vulnerability

Please report vulnerabilities privately through GitHub:
**[Report a vulnerability](https://github.com/kinglet-dev/laserfix/security/advisories/new)**
(Security tab → "Report a vulnerability").

Please don't open a public issue for security problems.

Include what you found, how to reproduce it (a sample SVG helps), and the
impact you expect.

## What laserfix trusts

Nothing in the SVG it repairs. laserfix reads one file or standard input,
and writes one new file or standard output. It never changes the file it
reads, never replaces an existing file, and writes through a temporary file
that is renamed into place, so a failure leaves no half-written file. It
uses no network and never follows references inside the file. Inputs over
its size and complexity limits are refused.

The machine-readable contact is published at
https://kinglet.dev/.well-known/security.txt (RFC 9116).
