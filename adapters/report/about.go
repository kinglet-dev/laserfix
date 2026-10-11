// Package report holds what every report format shows besides the changes.
package report

// About says which tool made a report, what it read and where it wrote.
type About struct {
	Version string // laserfix's version
	Input   string // the file read, or "standard input"
	Output  string // where the fixed copy went; "" when only previewing
}
