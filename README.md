# loom

Version: **1.2.0**

loom is a cross-platform Go CLI for UI layout analysis, generation planning,
translation, pattern catalog validation, transfer planning, and workflow
diagnostics.

The project ships as an analyzer, generator, and translator for moving UI layout
intent between SwiftUI, WinUI XAML, Qt, and JUCE. The current source includes
analyzer reports, transfer planning, component graphs, reviewable generator
scaffolds, target contracts, and manifest-directed project build bundles.

## Local developer installation

```sh
make build                       # compile and install ~/.local/bin/loom
make compile                     # compile only to .build/loom (CI/cross-builds)
make build PREFIX="$HOME/.local" # explicit installation prefix
```

`make install` is equivalent to `make build`. Put `$HOME/.local/bin` before
Homebrew on PATH. The installed executable is copied out of the checkout, so
moving the source repo does not break it. `scripts/build-local.sh` stages builds,
rejects cross-architecture installation, and retains previous installs under
`$HOME/.local/share/loom/installs/`. The current `install-info.txt` in that tool's
share directory records source path, commit, dirty state, Go version, and SHA-256.
Raw `go build` and release/CI scripts remain compile/package-only.
Run `python3 scripts/test-local-install.py` for isolated installer regression checks.

The installed pattern catalog is `$HOME/.local/share/loom/patterns`; verify it from outside the checkout with `loom patterns:validate --json`.


## What loom Does Today

- Parse WinUI XAML and normalize it into loom's shared layout model (`inspect:xaml`).
- Parse common SwiftUI layout/control constructs into the same shared model
  (`inspect:swiftui`, `inspect:source`).
- Parse common Qt QML, Qt Designer UI, and Qt C++ layout constructs into the
  same shared model (`inspect:qt`, `inspect:source`).
- Parse common JUCE C++ Component and widget constructs into the same shared
  model (`inspect:juce`, `inspect:source`).
- Extract intrinsic font material properties from supplied TrueType, OpenType,
  TrueType Collection, and WOFF font files or installed family names
  (`inspect:font`).
- Preserve WinUI Grid row/column definitions as layout metadata for transfer planning.
- Render parsed layouts as a compact plaintext ASCII tree (`inspect:ascii`).
- Validate and lint the `patterns` catalog (`patterns:validate`, `patterns:lint`).
- Transfer-plan layout compatibility in both WinUI → macOS and macOS/SwiftUI →
  Windows directions (`patterns:transfer`).
- Compare layout parity across supported source dialects (`inspect:parity`).
- Compare profile-normalized visual metrics and semantic visual roles such as
  toolbar, sidebar, media, tab, collection, emphasis, density, and native
  component boundaries (`inspect:visual-parity`).
- Audit accessibility/layout quality for unsupported boundaries, small targets,
  malformed or redundant structures, and scan-friendly risks (`accessibility:audit`).
- Run manifest validation (`config:validate` / `config:schema`).
- Discover source-tree component dependencies (`graph:components`).
- Run manifest-directed analysis builds that write validation, analysis, graph,
  generated scaffold, contract, transfer, parity, and summary artifacts
  (`project:build`).
- Generate reviewable WinUI XAML fragments, SwiftUI scaffolds, JUCE Component
  scaffolds, and target contracts (`generate:xaml`, `generate:swiftui`,
  `generate:juce`, `generate:contracts`).
- Report errors for source analysis, XAML parsing, manifests, and patterns
  (`inspect:errors`).
- Provide curated cross-platform error guidance and suggested fixes (`suggestions:os-errors`).
- Provide CLI self-diagnostics and guard reporting (`status`, `verify`,
  `checks:command-catalog`, `guards:summary`, `self-heal:plan`).
- Write deterministic LF output by default, with `--line-ending crlf` for Windows
  artifacts and `--line-ending native` for host-native text output.

## GCF interchange

loom 1.1.0 supports complete generic GCF snapshots alongside JSON, using pinned
gcf-go v1.8.0. Existing text/JSON defaults, generated native source and project
bundle artifact formats remain unchanged.

```sh
loom capabilities --format gcf
loom data encode report.json --format gcf
loom data decode report.gcf --format json
loom data stats report.json
loom inspect:source examples/sampleapp/contentview.swift --format gcf
loom inspect:xaml examples/sampleapp/mainwindow.xaml --format auto --output report.gcf
loom list --format gcf
```

`data` accepts a file or `-` for stdin and does not execute manifests. Encode
defaults to GCF; decode, stats and capabilities default to JSON. JSON report
commands accept `--format gcf|auto`; auto selects the smaller encoding, JSON on
ties, and falls back to JSON for unrepresentable values. Output paths are honored
exactly; choose an appropriate suffix and inspect the header for Auto output.
Normal overwrite and input-protection guards remain active.

Project manifests and visual profiles can be JSON or generic GCF; existing schema
and field validation still apply. Structured input/decoded manifests/profiles are
capped at 2 MiB, conversion and report output at 64 MiB. Invalid UTF-8, trailing
values and graph/session-delta profiles are rejected. Pattern catalogs stay JSON;
project bundle artifacts stay JSON/native text, while the CLI summary can be GCF.
GCF preserves report data; it does not make generated scaffolds production-ready
or establish token/performance savings. See THIRD_PARTY_NOTICES.md.

## patterns catalog

loom ships with its own `patterns` catalog. These files are public, neutral,
OS-agnostic examples that describe common interface elements such as buttons,
text, grids, lists, scroll regions, split views, toggles, and text input.

When you run loom from this repository, it uses `./patterns` by default. When
installed with `make build`, the same catalog is copied to
`$HOME/.local/share/loom/patterns` so the installed `loom` command can validate,
lint, list, export, and transfer-plan against loom's own pattern definitions.

Generator output is intentionally conservative. It is designed to be reviewable,
low-risk, and useful for handoff workflows.

## Quick Start

```sh
make build
./examples/sampleapp/analyze-sample-app.sh --overwrite
loom project:build examples/sampleapp/loom.json --output-dir examples/sampleapp/generated/project-build --overwrite --json
```

You should see generated analysis JSON and generated artifact drafts in
`examples/sampleapp/generated/`.

## Build And Test

```sh
make build      # installs ~/.local/bin/loom and its pattern catalog
make test       # runs go test ./...
go vet ./...    # static correctness check
```

## Command Quickstart

```sh
loom list
loom help <command>
loom status --json
loom verify --json
loom checks:command-catalog --json
loom inspect:source contentview.swift --json
loom inspect:swiftui contentview.swift --format json
loom inspect:qt mainwindow.qml --format json
loom inspect:juce MainComponent.cpp --format json
loom inspect:font Inter.ttf --json
loom inspect:font --family "Segoe UI" --json
loom inspect:xaml mainwindow.xaml --format json
loom inspect:parity contentview.swift --target mainwindow.qml --from swiftui --to qt --json
loom inspect:visual-parity contentview.swift --target mainwindow.xaml --from swiftui --to winui3 --profile visual-profile.json --json
loom inspect:visual-parity contentview.swift --target mainwindow.xaml --source-font Inter.ttf --target-font-family "Segoe UI" --json
loom inspect:ascii mainwindow.xaml --output layout.txt
loom inspect:ascii mainwindow.xaml --output layout.txt --line-ending crlf
loom accessibility:audit mainwindow.xaml --format json --fail-on warning
loom graph:components examples/sampleapp --format dot --output component-graph.dot
loom generate:xaml contentview.swift --output generated.xaml
loom generate:swiftui mainwindow.xaml --view-name MainWindowScaffold --output MainWindowScaffold.swift
loom generate:juce contentview.swift --class-name MainComponent --output MainComponent.hpp
loom generate:contracts contentview.swift --target winui3 --json
loom patterns:lint
loom patterns:transfer mainwindow.xaml --from winui3 --to macos
loom patterns:transfer contentview.swift --from swiftui --to windows
loom patterns:transfer mainwindow.qml --from qt --to windows
loom patterns:transfer contentview.swift --from swiftui --to juce
loom project:build examples/sampleapp/loom.json --output-dir examples/sampleapp/generated/project-build --overwrite --json
```

## Example Workflow

```sh
./examples/sampleapp/analyze-sample-app.sh --overwrite
```

The sample produces analysis, audit, transfer, component graph, project build,
parity, and visual-parity reports under `examples/sampleapp/generated/`.

## Recommended Screenshots

- [docs/screenshot-targets.md](docs/screenshot-targets.md) lists capture-ready
  command outputs and optional homepage screenshots that explain loom usage quickly.

## Repository Layout

- `cmd/loom`: Go CLI entrypoint
- `internal/loom`: CLI runtime, parsers, catalog logic, and reports
- `patterns`: OS-agnostic canonical layout and control pattern metadata
- `examples/sampleapp`: neutral public sample XAML and analysis script
- `docs`: operating docs, commands, and guides
- `README`: project context and command surface

## Documentation

- [docs/commands.md](docs/commands.md)
- [docs/json-schemas.md](docs/json-schemas.md)
- [docs/ai-agents.md](docs/ai-agents.md)
- [docs/architecture.md](docs/architecture.md)
- [docs/roadmap.md](docs/roadmap.md)
- [docs/support-policy.md](docs/support-policy.md)
- [docs/deprecations.md](docs/deprecations.md)
- [docs/release-checklist.md](docs/release-checklist.md)
- [docs/release-evidence.md](docs/release-evidence.md)
- [docs/screenshot-targets.md](docs/screenshot-targets.md)
- [CONTRIBUTING.md](CONTRIBUTING.md)
- [SECURITY.md](SECURITY.md)
- [TESTING.md](TESTING.md)

## License

loom is open source under the 0BSD license. See [LICENSE](LICENSE).

## Release Checklist

1. Ensure `VERSION` and `internal/loom/catalog.go` version constants match.
2. Run `make test`, `go vet ./...`, and `go run ./cmd/loom verify --json`.
3. Run the sample workflow with `--overwrite`.
4. Review [TODO.md](TODO.md) before cutting the next version.
5. Run `git add`, commit, and push.
6. Tag the release, e.g.:
   `git tag -a vX.Y.Z -m "loom X.Y.Z"; git push --tags`.

## Optional companion tools

`loom tools doctor` checks installations; `tools plan` previews workflows and `tools run` collects local reports in a new private directory. Normal commands continue to work without other Starlit tools. AI feedback requires a separate explicit report/peer invocation. See [CLI integration](docs/TOOL_INTEGRATION.md) for recipes, limits and snapshot ownership.
