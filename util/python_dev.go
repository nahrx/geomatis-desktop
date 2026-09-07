//go:build dev

package util

import (
	"fmt"
	"os/exec"
)

// This file is only compiled in development (`wails dev`, which automatically
// builds with the "dev" tag). In development, use whatever Python is
// installed on the machine's PATH, so contributors don't need a python-embed
// setup just to run `wails dev`.

// pythonExecutable resolves the Python interpreter to run pypy.py with.
func pythonExecutable() (string, error) {
	if p, err := exec.LookPath("python"); err == nil {
		return p, nil
	}
	if p, err := exec.LookPath("python3"); err == nil {
		return p, nil
	}
	return "", fmt.Errorf("no Python interpreter found on PATH. Install Python 3 and make sure it is available as 'python' or 'python3', and that opencv-python and numpy are installed (pip install opencv-python numpy)")
}

// pythonWorkDir returns the directory the Python process should run in, so
// `import pypy` can find pypy.py. Empty means: don't override cmd.Dir, use
// the process's own working directory -- which is the project root (where
// pypy.py lives) when running via `wails dev`.
func pythonWorkDir() (string, error) {
	return "", nil
}
