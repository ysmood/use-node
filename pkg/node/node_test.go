package node

import (
	"net/http"
	"os"
	"sync"
	"testing"

	"github.com/ysmood/got"
)

func TestGetNodeList(t *testing.T) {
	g := got.T(t)

	ver := getRemoteNodeList()[0].Ver

	validateVer(g, ">= v19.0.0", ver)
}

func TestGetPackageJSON(t *testing.T) {
	g := got.T(t)

	dir := "tmp/sub/path"
	g.MkdirAll(0755, dir)

	g.E(os.Chdir(dir))

	p := findPackageJSON()

	g.PathExists(p)
}

func TestGetVersion(t *testing.T) {
	g := got.T(t)

	validateVer(g, ">= v19.0.0", getNodeInfo("").Ver)
}

func TestGetNode(t *testing.T) {
	g := got.T(t)

	GetNodePath("")
	p := GetNodePath("")

	g.Has(p, "use-node")
}

func TestRegistries(t *testing.T) {
	g := got.T(t)

	const size = "39754090"

	wg := sync.WaitGroup{}
	for _, u := range newNode("v19.0.0").URLs() {
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

			g.Desc("%s %#v", u, res.Header).Eq(res.Header.Get("Content-Length"), size)
			wg.Done()
		}()
	}
	wg.Wait()
}
