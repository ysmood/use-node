package bun

import (
	"net/http"
	"path/filepath"
	"testing"

	"github.com/ysmood/got"
)

func TestGetBunList(t *testing.T) {
	g := got.T(t)

	ver := getRemoteBunList()[0]

	validateVer(g, ">= v1.0.0", ver.Ver())
}

func TestEngineBunNotExists(t *testing.T) {
	g := got.T(t)

	dir := filepath.FromSlash("tmp/engine-bun-not-exists")

	g.WriteFile(filepath.Join(dir, "package.json"), "{}")
	g.Chdir(dir)

	err := g.Panic(func() {
		getBunInfo("")
	})

	g.Has(err, `"engines.bun" is not found`)
}

func TestGetBun(t *testing.T) {
	g := got.T(t)

	GetBunPath(g.Context(), "latest", nil)
	p := GetBunPath(g.Context(), "latest", nil)

	g.Has(p, "use-node")
}

func TestRegistries(t *testing.T) {
	g := got.T(t)

	for _, u := range Bun("v1.3.13").URLs() {
		req, err := http.NewRequest("HEAD", u, nil)
		g.E(err)

		req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/111.0.0.0 Safari/537.36")

		res, err := http.DefaultClient.Do(req)
		g.E(err)
		g.Cleanup(func() {
			_ = res.Body.Close()
		})

		g.Desc("%s %#v", u, res.Header).Eq(res.StatusCode/100, 2)
	}
}
