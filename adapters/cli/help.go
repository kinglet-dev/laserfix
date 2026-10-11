package cli

// help is laserfix's usage text.
const help = `laserfix repairs an SVG so the laser burns each line once: it removes
shapes that copy an earlier shape. It never changes the file you give it.

Usage:
  laserfix [flags] FILE          show what would change; write nothing
  laserfix [flags] FILE -o OUT   write the repaired copy to OUT
  laserfix [flags] - -o -        read standard input, write standard output

Flags:
  -o OUT           write the repaired copy to OUT, a new file (an existing
                   file is never replaced), or - for standard output; the
                   report then goes to standard error
  --same LENGTH    lines nearer than this are the same line (default 0.02mm)
  --json           write the report as JSON (kinglet.report/v1)
  --version        show the version
  -h, --help       show this help

Lengths need a unit: mm, cm, in, pt or px.

Exit status: 0 nothing left to fix, 1 changes are needed (preview) or the
changes made the design worse, 2 the file couldn't be read or written.

Examples:
  laserfix grid.svg
  laserfix grid.svg -o grid-fixed.svg
  laserfix --json grid.svg > changes.json
  cat grid.svg | laserfix - -o - | laserlint -
`
