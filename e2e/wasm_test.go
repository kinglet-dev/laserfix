package e2e_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// The browser runs laserfix built for WebAssembly (GOOS=js GOARCH=wasm),
// reading the design from standard input and writing the repaired copy to
// standard output. These specs build that copy and
// run it under Node.js with Go's own runner, and require it to behave
// exactly like the native binary: same output, same errors, same exit code.
var _ = Describe("The WebAssembly build", Ordered, func() {
	var wasm, node, runner string

	BeforeAll(func() {
		// Arrange (shared): build the WebAssembly copy with the same version
		// as the native binary, and find Node.js and Go's runner for it.
		var err error
		node, err = exec.LookPath("node")
		Expect(err).NotTo(HaveOccurred(), "Node.js is needed to run the WebAssembly build; install it and run the tests again")
		goroot, err := exec.Command("go", "env", "GOROOT").Output()
		Expect(err).NotTo(HaveOccurred())
		runner = filepath.Join(strings.TrimSpace(string(goroot)), "lib", "wasm", "wasm_exec_node.js")
		wasm = filepath.Join(GinkgoT().TempDir(), "laserfix.wasm")
		build := exec.Command("go", "build", "-ldflags", "-X main.version=0.0.0-test", "-o", wasm, mainPackage)
		build.Env = append(os.Environ(), "GOOS=js", "GOARCH=wasm")
		out, err := build.CombinedOutput()
		Expect(err).NotTo(HaveOccurred(), string(out))
	})

	// inBrowserBuild runs the WebAssembly copy. Go's runner passes its whole
	// environment to the program, which must stay small, so it gets only
	// what Node.js itself needs.
	inBrowserBuild := func(stdin string, args ...string) (stdout, stderr string, code int) {
		cmd := exec.Command(node, append([]string{runner, wasm}, args...)...)
		cmd.Env = []string{}
		if root := os.Getenv("SystemRoot"); root != "" {
			cmd.Env = append(cmd.Env, "SystemRoot="+root)
		}
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

	example := func() string {
		b, err := os.ReadFile(filepath.Join("..", "examples", "stacked.svg"))
		Expect(err).NotTo(HaveOccurred())
		return string(b)
	}

	DescribeTable("matches the native binary",
		func(stdin func() string, args ...string) {
			// Arrange
			in := stdin()

			// Act
			wantOut, wantErr, wantCode := laserfix(in, args...)
			gotOut, gotErr, gotCode := inBrowserBuild(in, args...)

			// Assert
			Expect(gotOut).To(Equal(wantOut))
			Expect(gotErr).To(Equal(wantErr))
			Expect(gotCode).To(Equal(wantCode))
		},
		Entry("previewing a design with nothing to fix", func() string { return svg(square) }, "-"),
		Entry("previewing a copy pasted twice", func() string { return svg(square + square) }, "-"),
		Entry("previewing the stacked example", example, "-"),
		Entry("previewing the stacked example as JSON", example, "--json", "-"),
		Entry("writing the repaired copy", example, "-", "-o", "-"),
		Entry("writing the repaired copy with a JSON report", example, "--json", "-", "-o", "-"),
		Entry("with settings changed", example, "--same", "0.05mm", "--join", "0.1mm", "--report-gaps", "2mm", "-", "-o", "-"),
		Entry("on a file it can't read", func() string { return "<svg" }, "-"),
		Entry("on a bad flag", func() string { return svg(square) }, "--join", "far", "-"),
		Entry("showing its version", func() string { return "" }, "--version"),
	)
})
