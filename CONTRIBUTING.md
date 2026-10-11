# Contributing

Thanks for your interest in laserfix.

## Before you start

- **Open an issue first** for anything beyond a typo fix, so we can agree on
  the change before you spend time on it. A sample SVG that shows the problem
  helps a lot.
- **Security problems:** don't open an issue; follow [SECURITY.md](SECURITY.md).

## Contributor License Agreement

Outside contributions can only be merged after the contributor signs a
Contributor License Agreement (CLA). The CLA process is not set up yet, so pull
requests from outside contributors can't be merged for now. Issues, bug reports
and sample files are very welcome in the meantime.

## How changes are made

- Work on a branch (e.g. `feature/<short-name>`), never directly on `main`.
- Test first: write a failing Ginkgo test, then the code that makes it pass.
  Coverage must stay at 100% (`go tool ginkgo -r --cover`).
- Keep each commit to one change with its tests; every commit builds and passes.
- Commit messages: a short imperative subject (72 characters at most), a blank
  line, then why the change was made.
- No new dependencies without discussion; any dependency must use a license
  that allows unrestricted commercial use (MIT, Apache 2.0, BSD, ISC).
- Don't commit secrets, personal data, local file paths, build output or
  coverage files.
