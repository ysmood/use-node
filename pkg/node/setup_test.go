package node

import (
	"github.com/Masterminds/semver/v3"
	"github.com/ysmood/got"
)

func validateVer(g got.G, constrain string, ver *semver.Version) bool {
	c, err := semver.NewConstraint(constrain)
	g.E(err)

	return c.Check(ver)
}
