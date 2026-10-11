package output_test

import (
	"errors"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/kinglet-dev/laserfix/adapters/output"
)

var _ = Describe("Write cleans up when a step fails", func() {
	broken := errors.New("disk full")

	DescribeTable("returns the error and leaves the folder empty",
		func(step output.Step, removals int) {
			// Arrange
			dir := GinkgoT().TempDir()
			write, removed := output.Fail(step, broken)

			// Act
			err := write(filepath.Join(dir, "fixed.svg"), []byte("<svg/>"))

			// Assert
			Expect(err).To(MatchError(broken))
			entries, _ := os.ReadDir(dir)
			Expect(entries).To(BeEmpty())
			Expect(*removed).To(HaveLen(removals))
		},
		Entry("when the temporary file can't be made", output.Step("temp"), 1),
		Entry("when writing fails", output.Step("write"), 2),
		Entry("when flushing to disk fails", output.Step("sync"), 2),
		Entry("when setting permissions fails", output.Step("chmod"), 2),
		Entry("when closing fails", output.Step("close"), 2),
		Entry("when the rename fails", output.Step("rename"), 2),
	)

	It("drops the temporary file's name from a failed rename", func() {
		// Arrange
		dir := GinkgoT().TempDir()
		write, _ := output.Fail("rename", &os.LinkError{Op: "rename", Old: "a", New: "b", Err: broken})

		// Act
		err := write(filepath.Join(dir, "fixed.svg"), []byte("<svg/>"))

		// Assert
		Expect(err).To(Equal(broken))
	})
})
