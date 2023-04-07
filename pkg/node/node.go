package node

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/Masterminds/semver/v3"
	"github.com/ysmood/fetchup"
	"github.com/ysmood/gson"
	"github.com/ysmood/use-node/pkg/utils"
)

var cacheDir = filepath.Join(fetchup.CacheDir(), "use-node")
var versionsFile = filepath.Join(cacheDir, "versions.json")

func GetNodePath(required string, logger fetchup.Logger) string {
	utils.E(os.MkdirAll(cacheDir, 0755))

	n := getNodeInfo(required)

	nodePath := filepath.Join(cacheDir, n.String())

	if binExist(nodePath) {
		return nodePath
	}

	os.RemoveAll(nodePath)

	fu := fetchup.New(nodePath, n.URLs()...)

	if logger == nil {
		fu.Logger = fetchup.LoggerQuiet
	} else {
		fu.Logger = logger
	}

	utils.E(fu.Fetch())

	utils.E(fetchup.StripFirstDir(nodePath))

	return nodePath
}

type Node string

func (n Node) Ver() *semver.Version {
	v, err := semver.NewVersion(string(n))
	utils.E(err)
	return v
}

func (n Node) String() string {
	return string(n)
}

func getLocalNodeList() []Node {
	out := []Node{}

	list, err := os.ReadDir(cacheDir)
	utils.E(err)

	for _, d := range list {
		if !d.IsDir() || d.Name()[0] != 'v' {
			continue
		}

		out = append(out, Node(d.Name()))
	}

	return out
}

func getRemoteNodeList() []Node {
	us := []string{}
	for _, u := range famousRegistries {
		us = append(us, u+"/index.json")
	}
	fu := fetchup.New(versionsFile, us...)
	fu.SpeedPacketSize = 3 * 1024
	utils.E(fu.Fetch())

	return parseLocalNodeList()
}

func parseLocalNodeList() []Node {
	b, err := os.ReadFile(versionsFile)
	utils.E(err)

	list := gson.New(b)

	out := make([]Node, len(list.Arr()))

	for i, it := range list.Arr() {
		out[i] = Node(it.Get("version").Str())
	}

	return out
}

func getNodeInfo(required string) Node {
	if required == "" {
		p := findPackageJSON()
		if p == "" {
			panic("package.json not found")
		}

		b, err := os.ReadFile(p)
		utils.E(err)

		pkg := gson.New(b)

		required = pkg.Get("engines.node").Str()
		if required == "" {
			panic("Node version not found in package.json")
		}
	}

	if required == "latest" {
		return getRemoteNodeList()[0]
	}

	c, err := semver.NewConstraint(required)
	utils.E(err)

	for _, n := range getLocalNodeList() {
		if c.Check(n.Ver()) {
			return n
		}
	}

	for _, n := range getRemoteNodeList() {
		if c.Check(n.Ver()) {
			return n
		}
	}

	panic("No node version satisfies the requirement: " + required)
}

// recursively search for package.json
func findPackageJSON() string {
	d, err := os.Getwd()
	utils.E(err)

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

func BinPath(nodePath string) string {
	binPath := nodePath
	if runtime.GOOS != "windows" {
		binPath = filepath.Join(nodePath, "bin")
	}
	return binPath
}

func binExist(p string) bool {
	_, err := exec.Command(filepath.Join(BinPath(p), "node"), "-v").CombinedOutput()
	return err == nil
}
