# aless (another less)

[![CI](https://github.com/incu6us/another-less/actions/workflows/ci.yml/badge.svg)](https://github.com/incu6us/another-less/actions/workflows/ci.yml)
[![Release](https://img.shields.io/github/v/release/incu6us/another-less)](https://github.com/incu6us/another-less/releases/latest)
[![Go Report Card](https://goreportcard.com/badge/github.com/incu6us/another-less)](https://goreportcard.com/report/github.com/incu6us/another-less)
[![Go Reference](https://pkg.go.dev/badge/github.com/incu6us/another-less.svg)](https://pkg.go.dev/github.com/incu6us/another-less)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A TUI pager for exploring and analyzing tabular data -- excel for your logs.

Pipe in any data (CSV, TSV, JSON, logs, command output) and aless infers the structure, letting you filter, sort, pivot, search, and reshape interactively with vi-like keybindings.

Inspired by [nothing-less](https://github.com/mpryor/nothing-less), rewritten in Go for a single static binary with better performance.

## Install

### Homebrew

```bash
brew install incu6us/tap/aless
```

### Go

```bash
go install github.com/incu6us/another-less/cmd/aless@latest
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
