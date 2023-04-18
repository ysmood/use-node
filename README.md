# Overview

![usage](https://user-images.githubusercontent.com/1415488/229589650-c4112b86-f533-4205-a48f-44e88ebbc214.svg)

No need to setup anything, just run the `use-node`, then you can use the right node version under your shell.

- Works the same across Windows, macOS, and Linux.
- No options, no configs, just run `use-node` under your node project.
- Auto-detect the standard [engines](https://docs.npmjs.com/cli/v9/configuring-npm/package-json#engines) config recursively from all parent folders.
- Auto-choose the fastest registry to download node.
- Auto-cache the node binary.
- Lightweight and fast.
- Reports every side-effect it makes on the system, such as where the node is installed.

## Installation

Go to the [release page](https://github.com/ysmood/use-node/releases) and download the use-node binary, run the command below to install use-node to your [PATH]:

```bash
use-node -i
```

On macOS you need to do this before you can run the `use-node` binary:

```bash
xattr -r -d com.apple.quarantine /path/to/use-node
```

If you have golang installed:

```bash
go install github.com/ysmood/use-node@latest
```

## Use it outside node project

If you want to use it without the `package.json` file, just specify the node version you want to:

```bash
use-node v19
```

## Shell Setup

On unix like OS, add this line to your `.bashrc` or `.zshrc` files:

```bash
export PATH=$PATH:$(use-node -p v19.8.1)
```

If you want to run `use-node` whenever you [cd](<https://en.wikipedia.org/wiki/Cd_(command)>) to a directory, add this line to your `.bashrc` or `.zshrc` files:

```bash
eval "$(use-node -s)"
```

On Windows, run `use-node -p v19.8.1`, copy the output directory path and add it to [PATH].

## How it works

When you run the `use-node`. It first recognize the `engines` field in `package.json` and automatically download or use the best node version from cache.
Then it spawns a sub-shell and inject the node path to your PATH, your current stdio will transparently pipe to the sub-shell.

[path]: https://en.wikipedia.org/wiki/PATH_(variable)
