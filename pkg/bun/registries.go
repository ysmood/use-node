package bun

import (
	"fmt"
	"runtime"

	"github.com/ysmood/use-node/pkg/utils"
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

	suffix := ""
	if sys == "linux" && utils.IsMusl() {
		suffix = "-musl"
	}

	return []string{
		fmt.Sprintf("%s/bun-%s/bun-%s-%s%s.zip", downloadBase, b.String(), sys, arch, suffix),
	}
}
