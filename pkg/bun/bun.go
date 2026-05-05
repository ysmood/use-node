package bun

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.com/ysmood/fetchup"
	"github.com/ysmood/gson"
	"github.com/ysmood/use-node/pkg/utils"
)

var cacheDir = filepath.Join(fetchup.CacheDir(), "use-node", "bun")
var versionsFile = filepath.Join(cacheDir, "versions.json")

func GetBunPath(ctx context.Context, required string, logger fetchup.Logger) string {
	utils.E(os.MkdirAll(cacheDir, 0755))

	b := getBunInfo(required)

	bunPath := filepath.Join(cacheDir, b.String())

	if binExist(bunPath) {
		return bunPath
	}

	_ = os.RemoveAll(bunPath)

	fu := fetchup.New(bunPath, b.URLs()...)
	fu.Ctx = ctx

	if logger == nil {
		fu.Logger = fetchup.LoggerQuiet
	} else {
		fu.Logger = logger
	}

	utils.E(fu.Fetch())

	utils.E(fetchup.StripFirstDir(bunPath))

	utils.E(utils.WriteSentinel(bunPath, filepath.Join(BinPath(bunPath), bunBinName())))

	return bunPath
}

func bunBinName() string {
	if runtime.GOOS == "windows" {
		return "bun.exe"
	}
	return "bun"
}

type Bun string

func (b Bun) Ver() *semver.Version {
	v, err := semver.NewVersion(string(b))
	utils.E(err)
	return v
}

func (b Bun) String() string {
	return string(b)
}

func getLocalBunList() []Bun {
	out := []Bun{}

	list, err := os.ReadDir(cacheDir)
	utils.E(err)

	for _, d := range list {
		if !d.IsDir() || d.Name()[0] != 'v' {
			continue
		}

		out = append(out, Bun(d.Name()))
	}

	return out
}

func getRemoteBunList() []Bun {
	fu := fetchup.New(versionsFile, releasesAPI)
	utils.E(fu.Fetch())

	return parseLocalBunList()
}

func parseLocalBunList() []Bun {
	b, err := os.ReadFile(versionsFile)
	utils.E(err)

	list := gson.New(b)

	out := make([]Bun, len(list.Arr()))

	for i, it := range list.Arr() {
		out[i] = Bun(strings.TrimPrefix(it.Get("tag_name").Str(), "bun-"))
	}

	return out
}

func getBunInfo(required string) Bun {
	if required == "" {
		p := utils.FindPackageJSON()
		if p == "" {
			panic("package.json not found")
		}

		b, err := os.ReadFile(p)
		utils.E(err)

		pkg := gson.New(b)

		if pkg.Has("engines.bun") {
			required = pkg.Get("engines.bun").Str()
		} else {
			panic(`"engines.bun" is not found in: ` + p + ", for details: https://docs.npmjs.com/cli/v9/configuring-npm/package-json#engines")
		}
	}

	if required == "latest" {
		return getRemoteBunList()[0]
	}

	c, err := semver.NewConstraint(required)
	utils.E(err)

	for _, b := range getLocalBunList() {
		if c.Check(b.Ver()) {
			return b
		}
	}

	for _, b := range getRemoteBunList() {
		if c.Check(b.Ver()) {
			return b
		}
	}

	panic("No bun version satisfies the requirement: " + required)
}

// BinPath returns the directory containing the bun binary. Bun ships a single
// flat binary at the top of the extracted archive, so this is just bunPath.
func BinPath(bunPath string) string {
	return bunPath
}

func binExist(p string) bool {
	bin := filepath.Join(BinPath(p), bunBinName())
	if utils.CheckSentinel(p, bin) {
		return true
	}
	// Lazy migration for caches created before the sentinel existed:
	// validate by execing once, then record the fingerprint.
	if _, err := exec.Command(bin, "-v").CombinedOutput(); err != nil {
		return false
	}
	_ = utils.WriteSentinel(p, bin)
	return true
}
