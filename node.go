package main

import (
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/Masterminds/semver/v3"
	"github.com/ysmood/fetchup"
	"github.com/ysmood/gson"
)

var cacheDir = filepath.Join(fetchup.CacheDir(), "use-node")

func getNodePath() string {
	n := getNodeInfo()
	nodePath := filepath.Join(cacheDir, n.Ver.Original())

	if binExist(nodePath) {
		return nodePath
	}

	os.RemoveAll(nodePath)

	fu := fetchup.New(nodePath, n.URLs()...)

	E(fu.Fetch())

	E(fetchup.StripFirstDir(nodePath))

	return nodePath
}

type Node struct {
	Ver *semver.Version
}

func newNode(version string) Node {
	ver, err := semver.NewVersion(version)
	E(err)

	return Node{
		Ver: ver,
	}
}

func getLocalNodeList() []Node {
	out := []Node{}

	list, err := os.ReadDir(cacheDir)
	E(err)

	for _, d := range list {
		if !d.IsDir() || d.Name()[0] != 'v' {
			continue
		}

		out = append(out, newNode(d.Name()))
	}

	return out
}

func getRemoteNodeList() []Node {
	res, err := http.Get("https://nodejs.org/dist/index.json")
	E(err)
	defer func() {
		_ = res.Body.Close()
	}()

	list := gson.New(res.Body)

	out := make([]Node, len(list.Arr()))

	for i, it := range list.Arr() {
		out[i] = newNode(it.Get("version").Str())
	}

	return out
}

func getNodeInfo() Node {
	p := findPackageJSON()
	if p == "" {
		panic("package.json not found")
	}

	b, err := os.ReadFile(p)
	E(err)

	pkg := gson.New(b)

	required := pkg.Get("engines.node").Str()
	if required == "" {
		panic("node version not found in package.json")
	}

	c, err := semver.NewConstraint(required)
	E(err)

	for _, n := range getLocalNodeList() {
		if c.Check(n.Ver) {
			return n
		}
	}

	for _, n := range getRemoteNodeList() {
		if c.Check(n.Ver) {
			return n
		}
	}

	panic("no node version satisfies the requirement: " + required)
}

// recursively search for package.json
func findPackageJSON() string {
	d, err := os.Getwd()
	E(err)

	prev := ""

	for d != prev {
		p := filepath.Join(d, "package.json")
		if _, err := os.Stat(p); err == nil {
			return p
		}

		if d == "/" {
			break
		}

		prev = d
		d = filepath.Dir(d)
	}

	return ""
}

func binExist(p string) bool {
	bin := filepath.Join(p, "bin", "node")
	_, err := exec.Command(bin, "-v").CombinedOutput()
	return err == nil
}
