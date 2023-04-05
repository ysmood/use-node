package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
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
		p("  If the [node-version] is not specified, it will start a new shell with the version defined in the 'package.json'. If the version doesn't exist, it will be auto downloaded.")
		p("  For more doc: https://github.com/ysmood/use-node")
		p("")
		p("Examples:")
		p("")
		p("  use-node latest")
		p("  use-node v19.8.1")
		p("  use-node v17")
		p("")
		flag.PrintDefaults()
	}

	onlyPrint := flag.Bool("p", false, "Only print the node bin folder path outside use-node context")
	install := flag.Bool("i", false, "Install the use-node binary to one of the folders in PATH")

	flag.Parse()

	if *install {
		installSelfToPATH()
		return
	}

	ver := flag.Arg(0)

	nodePath := node.GetNodePath(ver)
	binPath := node.BinPath(nodePath)

	if *onlyPrint {
		if !isInUseNodeContext() {
			p(binPath)
		}
		return
	}

	if isInUseNodeContext() {
		p("Already in use-node context, please run exit before use-node again:", nodePath)
		return
	}

	os.Setenv(USE_NODE_SHELL, "true")
	os.Setenv(PATH, strings.Join([]string{binPath, getEnvWithoutOtherUseNode()}, string(os.PathListSeparator)))

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

func getEnvWithoutOtherUseNode() string {
	reg := regexp.MustCompile(fmt.Sprintf(`[^%v]+use-node[^%v]+%v?`, os.PathListSeparator, os.PathListSeparator, os.PathListSeparator))
	return reg.ReplaceAllString(cleanPath(), "")
}

func cleanPath() string {
	list := getPathList()

	m := map[string]struct{}{}

	for _, i := range list {
		m[i] = struct{}{}
	}

	list = []string{}
	for i := range m {
		list = append(list, i)
	}
	return strings.Join(list, string(os.PathListSeparator))
}

func isInUseNodeContext() bool {
	_, has := os.LookupEnv(USE_NODE_SHELL)
	return has
}

func getPathList() []string {
	return strings.Split(os.Getenv(PATH), string(os.PathListSeparator))
}

func installSelfToPATH() {
	path, err := os.Executable()
	utils.E(err)

	list := getPathList()

	for _, dir := range list {
		to := filepath.Join(dir, "use-node")
		if _, err := os.Stat(to); err == nil {
			p("Already installed:", to)
			return
		}
	}

	sort.Slice(list, func(i, j int) bool {
		return len(list[i]) < len(list[j])
	})

	for _, dir := range list {
		to := filepath.Join(dir, "use-node")

		f, err := os.Open(path)
		utils.E(err)

		info, err := f.Stat()
		utils.E(err)

		n, err := os.OpenFile(to, os.O_RDWR|os.O_CREATE|os.O_TRUNC, info.Mode())
		if err != nil {
			continue
		}

		_, err = io.Copy(n, f)
		if err == nil {
			p("Installed use-node to:", to)
			return
		}
	}

	panic("Failed to install use-node, no folder in PATH is writable")
}
