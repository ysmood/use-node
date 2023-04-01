# Overview

No need to setup anything, just run the `use-node`, then you can use the right node version under your shell.

## How it works

When you run the `use-node`. It first recognize the `engines` field in `package.json` and automatically download or use the best node version from cache.
Then it spawns a sub-shell and inject the node path to your PATH, your current stdio will transparently proxy to the sub-shell.
