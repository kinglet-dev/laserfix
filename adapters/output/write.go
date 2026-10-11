// Package output writes the repaired document to a new file without ever
// replacing an existing one, and without leaving a half-written file if
// anything fails.
//
// It reserves the name by creating it exclusively, writes the content to a
// temporary file beside it, then renames that over the reservation, which
// it owns (the write-temp-then-rename pattern for atomic replacement).
package output

import (
	"errors"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

// ErrExists means the output file is already there; it is never replaced.
var ErrExists = errors.New("it already exists: choose another name, or move the old file away first")

// file is the part of a temporary file Write uses.
type file interface {
	io.WriteCloser
	Name() string
	Sync() error
	Chmod(fs.FileMode) error
}

// system is the part of the file system Write uses, so tests can make
// each step fail (a port at the operating system's boundary).
type system struct {
	reserve func(name string) error
	temp    func(dir, pattern string) (file, error)
	rename  func(from, to string) error
	remove  func(name string) error
}

// osSystem is the real file system.
var osSystem = system{
	reserve: func(name string) error {
		f, err := os.OpenFile(name, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o644)
		if err == nil {
			// Nothing was written to it, and the rename replaces it, so a
			// failure to close it has no effect on the result.
			f.Close()
		}
		return err
	},
	temp:   func(dir, pattern string) (file, error) { return os.CreateTemp(dir, pattern) },
	rename: os.Rename,
	remove: os.Remove,
}

// Write writes data to a new file called name.
func Write(name string, data []byte) error { return osSystem.write(name, data) }

func (s system) write(name string, data []byte) error {
	if err := s.reserve(name); err != nil {
		if errors.Is(err, fs.ErrExist) {
			return ErrExists
		}
		return unwrapPath(err)
	}
	f, err := s.temp(filepath.Dir(name), "."+filepath.Base(name)+".*.tmp")
	if err != nil {
		s.remove(name)
		return unwrapPath(err)
	}
	err = fill(f, data)
	if err == nil {
		err = s.rename(f.Name(), name)
	}
	if err != nil {
		s.remove(f.Name())
		s.remove(name)
		return unwrapPath(err)
	}
	return nil
}

// fill writes data, flushes it to disk, gives the file the usual
// permissions for a new document, and closes it.
func fill(f file, data []byte) error {
	_, err := f.Write(data)
	if err == nil {
		err = f.Sync()
	}
	if err == nil {
		err = f.Chmod(0o644)
	}
	if cerr := f.Close(); err == nil {
		err = cerr
	}
	return err
}

// unwrapPath drops the paths from a file system error: the caller shows
// the name already, and the temporary file's name means nothing to people.
func unwrapPath(err error) error {
	if pe := (*fs.PathError)(nil); errors.As(err, &pe) {
		return pe.Err
	}
	if le := (*os.LinkError)(nil); errors.As(err, &le) {
		return le.Err
	}
	return err
}
