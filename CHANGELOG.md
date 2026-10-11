# Changelog

All notable changes to laserfix are listed here. It follows
[Semantic Versioning](https://semver.org/): the command line, the exit codes
and the JSON report (`kinglet.report/v1`) don't break within a major version.

## 0.1.0 (unreleased)

### Added

- Finds shapes that copy an earlier shape (same paint, same line within
  `--same`, 0.02 mm by default, wherever it starts and whichever way it
  runs) and removes the later copy.
- Shows what would change by default; `-o FILE` writes a repaired copy to
  a new file and never replaces one; `-` and `-o -` read standard input and
  write standard output.
- Measures line burned twice and crossings before and after, and exits 1
  if the copy came out worse.
- Text report for people and JSON report (`kinglet.report/v1`) with a
  `changes` list.
