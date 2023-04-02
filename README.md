# Overview

No need to setup anything, just run the `use-node`, then you can use the right node version under your shell.

- No options, no configs, just run `use-node` under your node project.
- Auto-detect the standard [engines](https://docs.npmjs.com/cli/v9/configuring-npm/package-json#engines) config recursively.
- Auto-choose the fastest registry to download node.
- Auto-cache the node binary.

## How it works

When you run the `use-node`. It first recognize the `engines` field in `package.json` and automatically download or use the best node version from cache.
Then it spawns a sub-shell and inject the node path to your PATH, your current stdio will transparently proxy to the sub-shell.
