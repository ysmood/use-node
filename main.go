package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/ysmood/use-node/pkg/node"
	"github.com/ysmood/use-node/pkg/utils"
)

const (
	PATH           = "PATH"
	USE_NODE_SHELL = "USE_NODE_SHELL"
)

func main() {
	flag.Usage = func() {
		p("Usage: use-node [node-version]")
		p("")
		p("  If the [node-version] is specified it only prints the local node path for the specified version. If the version doesn't exist, it will be auto downloaded.")
		p("  If the [node-version] is not specified, it will start a new shell with the version defined in the package.json .")
		p("")
		p("Example:")
		p("")
		p("  use-node v19.8.1")
		p("")
		flag.PrintDefaults()
	}
	flag.Parse()

	ver := flag.Arg(0)

	if _, has := os.LookupEnv(USE_NODE_SHELL); has {
		return
	}

	nodePath := node.GetNodePath(ver)
	binPath := node.BinPath(nodePath)

	if ver != "" {
		p(binPath)
		return
	}

	os.Setenv(USE_NODE_SHELL, "true")
	os.Setenv(PATH, strings.Join([]string{binPath, os.Getenv(PATH)}, string(os.PathListSeparator)))

	bin, err := Shell()
	utils.E(err)

	p("use-node:", nodePath)

	cmd := exec.Command(bin)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	_ = cmd.Run()
	os.Exit(cmd.ProcessState.ExitCode())
}

func p(v ...interface{}) {
	fmt.Println(v...)
}
