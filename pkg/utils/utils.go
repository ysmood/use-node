package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

const sentinelFile = ".use-node-ready"

// CheckSentinel reports whether binPath exists and matches the size+mtime
// fingerprint recorded by WriteSentinel in versionDir. A match means the
// install pipeline ran to completion and the binary hasn't been truncated
// or replaced since.
func CheckSentinel(versionDir, binPath string) bool {
	info, err := os.Stat(binPath)
	if err != nil {
		return false
	}
	data, err := os.ReadFile(filepath.Join(versionDir, sentinelFile))
	if err != nil {
		return false
	}
	var size, mtime int64
	if _, err := fmt.Sscanf(string(data), "%d %d", &size, &mtime); err != nil {
		return false
	}
	return size == info.Size() && mtime == info.ModTime().UnixNano()
}

// WriteSentinel records the size and mtime of binPath in a sentinel file
// inside versionDir.
func WriteSentinel(versionDir, binPath string) error {
	info, err := os.Stat(binPath)
	if err != nil {
		return err
	}
	content := fmt.Sprintf("%d %d", info.Size(), info.ModTime().UnixNano())
	return os.WriteFile(filepath.Join(versionDir, sentinelFile), []byte(content), 0644)
}

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
