package utils

import (
	"os"
	"path/filepath"
)

func E(err error) {
	if err != nil {
		panic(err)
	}
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
