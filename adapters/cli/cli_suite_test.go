package cli_test

import (
	"bytes"
	"errors"
	"io"
	"io/fs"
	"strings"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/kinglet-dev/laserfix/adapters/cli"
	"github.com/kinglet-dev/laserfix/adapters/output"
	"github.com/kinglet-dev/laserfix/adapters/svg"
	"github.com/kinglet-dev/laserfix/app"
)

func TestCLI(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "CLI Suite")
}

// doc is a 100 mm square document holding the given markup.
func doc(body string) string {
	return `<svg xmlns="http://www.w3.org/2000/svg" width="100mm" height="100mm" viewBox="0 0 100 100">` +
		body + `</svg>`
}

// sq is a 10 mm black hairline square with its corner at (10, 10).
const sq = `<rect x="10" y="10" width="10" height="10" fill="none" stroke="#000"/>`

// world is a run's surroundings in memory: files to read, files written,
// and what went to standard input, output and error.
type world struct {
	files   map[string]string
	written map[string]string
	opened  int // files opened and not yet closed
	stdin   string
	stdout  bytes.Buffer
	stderr  bytes.Buffer
}

func newWorld() *world {
	return &world{files: map[string]string{}, written: map[string]string{}}
}

// run runs laserfix in the world and returns the exit code.
func (w *world) run(args ...string) int { return w.runWith(func(*cli.Deps) {}, args...) }

// runWith runs laserfix after change has altered its surroundings.
func (w *world) runWith(change func(*cli.Deps), args ...string) int {
	d := cli.Deps{
		Stdin:   strings.NewReader(w.stdin),
		Stdout:  &w.stdout,
		Stderr:  &w.stderr,
		Open:    w.open,
		Write:   w.write,
		Fixer:   app.NewFixer(svg.Document{}, svg.Document{}),
		Version: "0.1.0",
	}
	change(&d)
	return cli.Run(args, d)
}

func (w *world) open(name string) (io.ReadCloser, error) {
	s, ok := w.files[name]
	if !ok {
		return nil, fs.ErrNotExist
	}
	w.opened++
	return closer{strings.NewReader(s), w}, nil
}

// closer is an open file that the world counts until it is closed.
type closer struct {
	io.Reader
	w *world
}

func (c closer) Close() error {
	c.w.opened--
	return nil
}

func (w *world) write(name string, data []byte) error {
	if _, ok := w.files[name]; ok {
		return output.ErrExists
	}
	if name == "broken.svg" {
		return errBroken
	}
	w.written[name] = string(data)
	return nil
}

var errBroken = errors.New("disk full")
