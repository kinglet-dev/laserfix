package output

import "io/fs"

// Step names a step of Write that a test can make fail.
type Step string

// Fail returns Write with one step failing with err, and the names it
// removed afterward, for tests of the cleanup.
func Fail(step Step, err error) (write func(name string, data []byte) error, removed *[]string) {
	s := osSystem
	var gone []string
	s.remove = func(name string) error {
		gone = append(gone, name)
		return osSystem.remove(name)
	}
	switch step {
	case "temp":
		s.temp = func(string, string) (file, error) { return nil, err }
	case "rename":
		s.rename = func(string, string) error { return err }
	default: // a step of the temporary file
		s.temp = func(dir, pattern string) (file, error) {
			f, e := osSystem.temp(dir, pattern)
			return failing{f, step, err}, e
		}
	}
	return s.write, &gone
}

// failing is a temporary file whose step fails.
type failing struct {
	file
	step Step
	err  error
}

func (f failing) Write(p []byte) (int, error) {
	return f.fail("write", func() (int, error) { return f.file.Write(p) })
}
func (f failing) Sync() error {
	_, err := f.fail("sync", func() (int, error) { return 0, f.file.Sync() })
	return err
}
func (f failing) Chmod(m fs.FileMode) error {
	_, err := f.fail("chmod", func() (int, error) { return 0, f.file.Chmod(m) })
	return err
}
func (f failing) Close() error {
	_, err := f.fail("close", func() (int, error) { return 0, f.file.Close() })
	return err
}

func (f failing) fail(step Step, real func() (int, error)) (int, error) {
	n, err := real()
	if f.step == step {
		return 0, f.err
	}
	return n, err
}
