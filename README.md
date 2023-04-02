# Overview

No need to setup anything, just run the `use-node`, then you can use the right node version under your shell.

- Auto-choose the fastest registry to download node
- Auto-cache the node binary
- Auto-detect the standard `engines` config in `package.json`

## How it works

When you run the `use-node`. It first recognize the `engines` field in `package.json` and automatically download or use the best node version from cache.
Then it spawns a sub-shell and inject the node path to your PATH, your current stdio will transparently proxy to the sub-shell.
