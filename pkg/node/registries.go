package node

import (
	"fmt"
	"runtime"
)

var famousRegistries = []string{
	"https://nodejs.org/dist",
	"https://mirrors.tuna.tsinghua.edu.cn/nodejs-release",
	"https://cdn.npmmirror.com/binaries/node",
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

	list := []string{}

	for _, r := range famousRegistries {
		list = append(list, fmt.Sprintf(
			"%s/%s/node-%s-%s-%s.%s",
			r,
			n.Ver.Original(),
			n.Ver.Original(),
			sys,
			arch,
			ext,
		))
	}

	return list
}
