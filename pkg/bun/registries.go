package bun

import (
	"fmt"
	"runtime"
)

const releasesAPI = "https://api.github.com/repos/oven-sh/bun/releases?per_page=100"

const downloadBase = "https://github.com/oven-sh/bun/releases/download"

func (b Bun) URLs() []string {
	sys := runtime.GOOS

	arch := runtime.GOARCH
	switch arch {
	case "amd64":
		arch = "x64"
	case "arm64":
		arch = "aarch64"
	}

	return []string{
		fmt.Sprintf("%s/bun-%s/bun-%s-%s.zip", downloadBase, b.String(), sys, arch),
	}
}
