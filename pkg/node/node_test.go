package node

import (
	"net/http"
	"path/filepath"
	"sync"
	"testing"

	"github.com/ysmood/got"
)

func TestGetNodeList(t *testing.T) {
	g := got.T(t)

	ver := getRemoteNodeList()[0]

	validateVer(g, ">= v19.0.0", ver.Ver())
}

func TestGetPackageJSON(t *testing.T) {
	g := got.T(t)

	dir := filepath.FromSlash("tmp/sub/path")
	g.MkdirAll(0755, dir)

	g.Chdir(dir)

	p := findPackageJSON()

	g.PathExists(p)
}

func TestEngineNodeNotExists(t *testing.T) {
	g := got.T(t)

	dir := filepath.FromSlash("tmp/engine-node-not-exists")

	g.WriteFile(filepath.Join(dir, "package.json"), "{}")
	g.Chdir(dir)

	err := g.Panic(func() {
		getNodeInfo("")
	})

	g.Has(err, `"engines.node" is not found`)
}

func TestGetVersion(t *testing.T) {
	g := got.T(t)

	validateVer(g, ">= v19.0.0", getNodeInfo("").Ver())
}

func TestGetNode(t *testing.T) {
	g := got.T(t)

	GetNodePath(g.Context(), "", nil)
	p := GetNodePath(g.Context(), "", nil)

	g.Has(p, "use-node")
}

func TestRegistries(t *testing.T) {
	g := got.T(t)

	size := ""

	wg := sync.WaitGroup{}
	for _, u := range Node("v19.0.0").URLs() {
		wg.Add(1)
		u := u
		go func() {
			req, err := http.NewRequest("", u, nil)
			g.E(err)

			req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/111.0.0.0 Safari/537.36")

			res, err := http.DefaultClient.Do(req)
			g.E(err)
			g.Cleanup(func() {
				_ = res.Body.Close()
			})

			actual := res.Header.Get("Content-Length")

			if size == "" {
				size = actual
			}

			g.Desc("%s %#v", u, res.Header).Eq(actual, size)
			wg.Done()
		}()
	}
	wg.Wait()
}
