package node

import (
	"log"
	"net/http"
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

func GetNodePath(required string) string {
	utils.E(os.MkdirAll(cacheDir, 0755))

	n := getNodeInfo(required)

	nodePath := filepath.Join(cacheDir, n.Ver.Original())

	if binExist(nodePath) {
		return nodePath
	}

	os.RemoveAll(nodePath)

	fu := fetchup.New(nodePath, n.URLs()...)
	fu.Logger = log.New(os.Stdout, "", 0)

	utils.E(fu.Fetch())

	utils.E(fetchup.StripFirstDir(nodePath))

	return nodePath
}

type Node struct {
	Ver *semver.Version
}

func newNode(version string) Node {
	ver, err := semver.NewVersion(version)
	utils.E(err)

	return Node{
		Ver: ver,
	}
}

func getLocalNodeList() []Node {
	out := []Node{}

	list, err := os.ReadDir(cacheDir)
	utils.E(err)

	for _, d := range list {
		if !d.IsDir() || d.Name()[0] != 'v' {
			continue
		}

		out = append(out, newNode(d.Name()))
	}

	return out
}

func getRemoteNodeList() []Node {
	us := []string{}
	for _, u := range famousRegistries {
		us = append(us, u+"/index.json")
	}
	fu := fetchup.New("", us...)
	fu.SpeedPacketSize = 3 * 1024

	res, err := http.Get(fu.FastestURL())
	utils.E(err)
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
			panic("node version not found in package.json")
		}
	}

	c, err := semver.NewConstraint(required)
	utils.E(err)

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
