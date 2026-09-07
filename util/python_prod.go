//go:build !dev

package util

import (
	"fmt"
	"os"
	"path/filepath"
)

// This file is compiled for everything except `wails dev` (i.e. `go build`
// and `wails build`, which do not set the "dev" tag). In a built app, use the
// Python runtime bundled next to the executable, so end users don't need
// Python installed on their machine.

// pythonExecutable resolves the Python interpreter to run pypy.py with:
// python-embed/python.exe next to the running executable.
func pythonExecutable() (string, error) {
	exePath, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("failed to resolve application path: %w", err)
	}
	return filepath.Join(filepath.Dir(exePath), "python-embed", "python.exe"), nil
}

// pythonWorkDir returns the directory the Python process should run in: the
// executable's own directory, so `import pypy` finds pypy.py inside
// python-embed/ (see python313._pth) regardless of the process's cwd.
func pythonWorkDir() (string, error) {
	exePath, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("failed to resolve application path: %w", err)
	}
	return filepath.Dir(exePath), nil
}
