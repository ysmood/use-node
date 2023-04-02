package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/ysmood/use-node/pkg/node"
	"github.com/ysmood/use-node/pkg/utils"
)

func main() {
	nodePath := node.GetNodePath()
	binPath := filepath.Join(nodePath, "bin")

	const osPATH = "PATH"

	os.Setenv(osPATH, strings.Join([]string{binPath, os.Getenv(osPATH)}, string(os.PathListSeparator)))

	bin, err := Shell()
	utils.E(err)

	cmd := exec.Command(bin)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	_ = cmd.Run()
	os.Exit(cmd.ProcessState.ExitCode())
}
