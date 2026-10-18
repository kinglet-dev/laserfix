# Changelog

All notable changes to laserfix are listed here. It follows
[Semantic Versioning](https://semver.org/): the command line, the exit codes
and the JSON report (`kinglet.report/v1`) don't break within a major version.

## 0.2.0 (2026-10-18)

### Added

- A WebAssembly copy of laserfix for web browsers, `laserfix_<version>.wasm`, published with each release alongside the `wasm_exec.js` that runs it and Go's license for that file. It reads the design from standard input and writes the repaired copy to standard output, like `laserfix - -o -`, with exactly the same report and exit status as the native program; a new end-to-end test checks this on every change.

## 0.1.0 (2026-10-11)

### Added

- Finds shapes that copy an earlier shape (same paint, same line within
  `--same`, 0.02 mm by default, wherever it starts and whichever way it
  runs) and removes the later copy.
- Shows what would change by default; `-o FILE` writes a repaired copy to
  a new file and never replaces one; `-` and `-o -` read standard input and
  write standard output.
- Trims line lying on earlier line (same stroke, within `--same`,
  running alongside it), another shape's or the shape's own, from the
  later line; the earlier one stays whole and curves stay curves.
- Repeats its fixes until a pass changes nothing (up to five passes), so
  one run fixes everything it can.
- Joins open pieces whose ends meet within `--join` (0.05 mm by default)
  into one path, closing it when the chain comes back to its start.
- Reports what it leaves as it is: junctions of three or more ends, gaps
  up to `--report-gaps` (1 mm by default), and line still burned twice by
  shapes of different colors or filled shapes.
- Measures line burned twice and crossings before and after, and exits 1
  if the copy came out worse or something was left as it is.
- Text report for people and JSON report (`kinglet.report/v1`) with a
  `changes` list and laserlint-style `findings`.
