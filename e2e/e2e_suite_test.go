// Package e2e runs the built laserfix binary like a user would. The binary
// is built with coverage on, so the composition root (cmd/laserfix), which
// unit tests can't reach, is measured too (engineering rules, section 3).
package e2e_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

const mainPackage = "github.com/kinglet-dev/laserfix/cmd/laserfix"

var binary, coverDir string

func TestE2E(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "End-to-End Suite")
}

var _ = BeforeSuite(func() {
	// Arrange (shared): build into a folder so Go picks the file name for
	// this OS, then find it.
	binDir := GinkgoT().TempDir()
	coverDir = GinkgoT().TempDir()
	build := exec.Command("go", "build", "-cover", "-ldflags", "-X main.version=0.0.0-test", "-o", binDir, mainPackage)
	out, err := build.CombinedOutput()
	Expect(err).NotTo(HaveOccurred(), string(out))
	found, _ := filepath.Glob(filepath.Join(binDir, "laserfix*"))
	Expect(found).To(HaveLen(1))
	binary = found[0]
})

var _ = AfterSuite(func() {
	// Assert: every statement of the composition root ran.
	out, err := exec.Command("go", "tool", "covdata", "percent", "-i", coverDir, "-pkg", mainPackage).CombinedOutput()
	Expect(err).NotTo(HaveOccurred(), string(out))
	Expect(string(out)).To(ContainSubstring("coverage: 100.0% of statements"))
})

// laserfix runs the binary with the given standard input and arguments.
func laserfix(stdin string, args ...string) (stdout, stderr string, code int) {
	cmd := exec.Command(binary, args...)
	cmd.Env = append(os.Environ(), "GOCOVERDIR="+coverDir)
	cmd.Stdin = stringsReader(stdin)
	var o, e buffer
	cmd.Stdout, cmd.Stderr = &o, &e
	err := cmd.Run()
	if exit, ok := err.(*exec.ExitError); ok {
		return o.String(), e.String(), exit.ExitCode()
	}
	Expect(err).NotTo(HaveOccurred())
	return o.String(), e.String(), 0
}
