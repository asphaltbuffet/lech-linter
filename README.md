# lech-linter

A linter that checks whether you spelled "Lechlitner" correctly.

That's it. That's the whole tool.

```
$ lech-linter .git/config
.git/config:2:8: possible misspelling "Lechniter" (did you mean "Lechlitner"?)
.git/config:7:2: possible misspelling "Lechlitnre" (did you mean "Lechlitner"?)
Error: 2 possible misspellings
```

---

## Usage

```
lech-linter [file ...] [flags]
```

Pass one or more files to check them. Pass nothing to read from stdin.

```bash
# Check a file
lech-linter CONTRIBUTORS.md

# Check several files
lech-linter docs/*.md

# Pipe something in
git log --oneline | lech-linter
```

Exit code is `0` if nothing is wrong, `1` if Lechlitner has been wronged.

Output format: `file:line:col: possible misspelling "X" (did you mean "Lechlitner"?)`

---

## GitHub Actions

Use lech-linter directly in your workflows to catch misspellings of "Lechlitner" in pull requests and pushes.

```yaml
- name: Run lech-linter
  uses: asphaltbuffet/lech-linter@v0
  with:
    files: '**/*.md'
```

### Inputs

| Input   | Required | Default | Description                                          |
|---------|----------|---------|------------------------------------------------------|
| `files` | No       | `''`    | Space-separated list of files or glob patterns to lint |

### Example workflow

```yaml
name: lech-linter

on: [push, pull_request]

jobs:
  lint:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Run lech-linter
        uses: asphaltbuffet/lech-linter@v0
        with:
          files: '**/*.md CONTRIBUTORS.txt'
```

---

## Installation

All packages include shell completions (bash, fish, zsh) and a man page. No extra steps required.

**Nix (flakes)**

```nix
# flake.nix inputs
lech-linter.url = "github:asphaltbuffet/lech-linter";
```

```nix
# home-manager
imports = [ inputs.lech-linter.homeManagerModules.default ];
```

Or try it without installing:

```bash
nix run github:asphaltbuffet/lech-linter -- --help
```

**Debian / Ubuntu**

```bash
dpkg -i lech-linter_*.deb
```

**RHEL / Fedora**

```bash
rpm -i lech-linter_*.rpm
```

**Alpine**

```bash
apk add --allow-untrusted lech-linter_*.apk
```

**Tarball (macOS, Linux, Windows)**

Download from the [releases page](https://github.com/asphaltbuffet/lech-linter/releases), extract, and put the binary somewhere on your `PATH`.

**Go**

```bash
go install github.com/asphaltbuffet/lech-linter@latest
```

Completions and man page not included. You get what you pay for.

---

## How it works

Any word between 7 and 13 characters with a [Levenshtein distance](https://en.wikipedia.org/wiki/Levenshtein_distance) ≤ 3 from "lechlitner" (case-insensitive) is flagged — unless it actually *is* "lechlitner", in which case well done.

The length filter exists so the tool doesn't flag "lech" or "supercalifragilistic" as possible misspellings. You're welcome.

Column numbers are byte offsets, not character offsets. If your source file contains multi-byte UTF-8 characters before a misspelling on the same line, the reported column will be slightly off. Given that the target name is ASCII and this tool exists to catch typos in people's names rather than parse Unicode poetry, this seems fine.

---

## Contributing

If you've found a word that is or isn't being flagged and you feel strongly about it, open an issue. The threshold and length filter are tunable if there's a compelling case.

This project uses [Jujutsu](https://github.com/jj-vcs/jj) for version control, [mise](https://mise.jdx.dev/) as the task runner, and [Nix flakes](https://nixos.wiki/wiki/Flakes) for reproducible builds.

```bash
mise run test   # run tests
mise run lint   # lint
mise run build  # build to ./dist/
```

---

*Named after Ben Lechlitner, who presumably got tired of seeing his name mangled.*
