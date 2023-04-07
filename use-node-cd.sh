__use-node-cd() {
    \cd "$@" || return $?
    if node=$(use-node -p 2>/dev/null); then
        if [[ ":$PATH:" != *":$node:"* ]]; then
            export PATH=$node:$PATH
        fi
    fi
}

alias cd=__use-node-cd
