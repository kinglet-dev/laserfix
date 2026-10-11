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
- Trims line lying on an earlier shape's line (same stroke, within
  `--same`, running alongside it) from the later shape; the earlier one
  stays whole and curves stay curves.
- Joins open pieces whose ends meet within `--join` (0.05 mm by default)
  into one path, closing it when the chain comes back to its start.
- Reports what it leaves as it is: junctions of three or more ends, gaps
  up to `--report-gaps` (1 mm by default), and line still burned twice by
  shapes of different colors or filled shapes.
- Measures line burned twice and crossings before and after, and exits 1
  if the copy came out worse or something was left as it is.
- Text report for people and JSON report (`kinglet.report/v1`) with a
  `changes` list and laserlint-style `findings`.
