# Overview

No need to setup anything, just run the `use-node`, then you can use the right node version under your shell.

- Works the same across Windows, macOS, and Linux.
- No options, no configs, just run `use-node` under your node project.
- Auto-detect the standard [engines](https://docs.npmjs.com/cli/v9/configuring-npm/package-json#engines) config recursively from all parent folders.
- Auto-choose the fastest registry to download node.
- Auto-cache the node binary.
- Reports every side-effect it makes on the system, such as where the node is installed.

## Installation

Go to the [release page](https://github.com/ysmood/use-node/releases) to download the use-node binary to one of your PATH.

If you have golang installed: `go install github.com/ysmood/use-node@latest`.

## How it works

When you run the `use-node`. It first recognize the `engines` field in `package.json` and automatically download or use the best node version from cache.
Then it spawns a sub-shell and inject the node path to your PATH, your current stdio will transparently pipe to the sub-shell.

## Use it outside node project

On unix like OS, add this line to your `.bashrc` or `.zshrc` files:

```bash
export PATH=$PATH:$(use-node v19.8.1)
```

On Windows, run `use-node v19.8.1`, copy the output then add it the [PATH config](https://www.java.com/en/download/help/path.html).
