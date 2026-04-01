# aless (another less)

A TUI pager for exploring and analyzing tabular data -- excel for your logs.

Pipe in any data (CSV, TSV, JSON, logs, command output) and aless infers the structure, letting you filter, sort, pivot, search, and reshape interactively with vi-like keybindings.

Inspired by [nothing-less](https://github.com/mpryor/nothing-less), rewritten in Go for a single static binary with better performance.

## Install

```bash
go install github.com/slavka/another-less/cmd/aless@latest
```

Or build from source:

```bash
make build
# binary at ./bin/aless
```

## Usage

```bash
# Pipe data
cat data.csv | aless
kubectl logs pod/myapp | aless
tail -f /var/log/syslog | aless

# Open files
aless access.log
aless data.csv
```

## Key Bindings (vim mode)

| Key | Action |
|-----|--------|
| `j` / `k` | Scroll down / up |
| `h` / `l` | Scroll left / right |
| `Ctrl-d` / `Ctrl-u` | Half page down / up |
| `Ctrl-f` / `Ctrl-b` | Full page down / up |
| `g` / `G` | Top / bottom |
| `0` / `$` | Line start / end |
| `/` | Search |
| `n` / `N` | Next / previous match |
| `q` | Quit |
| `Q` | Pipe buffer to stdout and quit |

## License

MIT
