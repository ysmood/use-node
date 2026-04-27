package node

import (
	"fmt"
	"runtime"

	"github.com/ysmood/use-node/pkg/utils"
)

var famousRegistries = []string{
	"https://nodejs.org/dist",
	"https://mirrors.tuna.tsinghua.edu.cn/nodejs-release",
	"https://cdn.npmmirror.com/binaries/node",
}

// muslRegistries hosts the community-built musl binaries. The official
// nodejs.org mirrors only ship glibc-linked builds.
var muslRegistries = []string{
	"https://unofficial-builds.nodejs.org/download/release",
}

// registries returns the URL bases to use for both the version index and
// per-version downloads, picking musl mirrors on a musl Linux system.
func registries() []string {
	if runtime.GOOS == "linux" && utils.IsMusl() {
		return muslRegistries
	}
	return famousRegistries
}

func (n Node) URLs() []string {
	ext := "tar.gz"

	sys := runtime.GOOS
	if sys == "windows" {
		sys = "win"
		ext = "zip"
	}

	arch := runtime.GOARCH
	if arch == "amd64" {
		arch = "x64"
	}

	suffix := ""
	if sys == "linux" && utils.IsMusl() {
		suffix = "-musl"
	}

	list := []string{}

	for _, r := range registries() {
		list = append(list, fmt.Sprintf(
			"%s/%s/node-%s-%s-%s%s.%s",
			r,
			n.String(),
			n.String(),
			sys,
			arch,
			suffix,
			ext,
		))
	}

	return list
}
