package utils

import (
	"os"
	"path/filepath"
	"runtime"
)

func E(err error) {
	if err != nil {
		panic(err)
	}
}

// IsMusl reports whether the current Linux system uses musl libc (e.g. Alpine).
// Always false on non-Linux. Detection is based on the presence of the musl
// dynamic loader at /lib/ld-musl-*.
func IsMusl() bool {
	if runtime.GOOS != "linux" {
		return false
	}
	matches, _ := filepath.Glob("/lib/ld-musl-*")
	return len(matches) > 0
}

// FindPackageJSON walks up from the current working directory looking for a package.json.
// Returns "" if none is found before reaching the filesystem root.
func FindPackageJSON() string {
	d, err := os.Getwd()
	E(err)

	prev := ""

	for d != prev {
		p := filepath.Join(d, "package.json")
		if _, err := os.Stat(p); err == nil {
			return p
		}

		prev = d
		d = filepath.Dir(d)
	}

	return ""
}
