# loom commands

## Command Surface (Curated by command family)

### Health & Guardrails

- `status`
- `verify`
- `checks:command-catalog`
- `guards:summary`
- `self-heal:plan`

### Validation & Configuration

- `config:validate`
- `config:schema`

### Accessibility

- `accessibility:audit`

### Inspection

- `inspect:source`
- `inspect:swiftui`
- `inspect:qt`
- `inspect:juce`
- `inspect:xaml`
- `inspect:ascii`
- `inspect:errors`
- `inspect:font`
- `inspect:parity`
- `inspect:visual-parity`
- `graph:components`

### Generation

- `generate:xaml`
- `generate:swiftui`
- `generate:juce`
- `generate:contracts`

### Patterns

- `patterns:list`
- `patterns:show`
- `patterns:validate`
- `patterns:lint`
- `patterns:export`
- `patterns:transfer`

### Project Packaging

- `project:build`

### Diagnostics

- `suggestions:os-errors`

## Generator And Translator Commands

Generator output is reviewable scaffold output. It preserves unsupported
component boundaries instead of silently inventing native behavior, and it keeps
file mutation behind explicit output or owned-region flags.

## Quick Recipes

- Start with health checks:
  `loom status --json` -> `loom verify --json` -> `loom checks:command-catalog --json`
- Inspect one file across formats:
  `loom inspect:swiftui examples/sampleapp/contentview.swift --format json`
  `loom inspect:xaml examples/sampleapp/mainwindow.xaml --format json`
- Generate and then handoff:
  `loom generate:xaml examples/sampleapp/contentview.swift --output /tmp/generated.xaml`
  `loom generate:juce examples/sampleapp/contentview.swift --class-name MainComponent --output /tmp/MainComponent.hpp`

## Notes

- `--quiet` suppresses successful write confirmations.
- `--verbose` prints write diagnostics to stderr.
- `--line-ending lf|crlf|native` controls text output line endings. loom uses
  `lf` by default for deterministic artifacts.
- `--overwrite` is required when a command with `--output` would replace an
  existing file.
- Most read/write report commands accept `--json`.
- Most commands accept `--help` for catalog-driven synopsis.

## examples

```sh
loom status --json
loom verify --json
loom checks:command-catalog --json
loom guards:summary
loom self-heal:plan
loom inspect:source contentview.swift --json
loom inspect:swiftui contentview.swift --format json
loom inspect:qt mainwindow.qml --format json
loom inspect:juce MainComponent.cpp --format json
loom inspect:font Inter.ttf --json
loom inspect:font --family "Segoe UI" --json
loom inspect:errors mainwindow.xaml --kind xaml --json --fail-on error
loom inspect:parity contentview.swift --target mainwindow.qml --from swiftui --to qt --json
loom inspect:visual-parity contentview.swift --target mainwindow.xaml --from swiftui --to winui3 --profile visual-profile.json --json
loom inspect:visual-parity contentview.swift --target mainwindow.xaml --source-font Inter.ttf --target-font-family "Segoe UI" --json
loom graph:components examples/sampleapp --format json
loom graph:components examples/sampleapp --format dot --output component-graph.dot
loom generate:xaml contentview.swift --output generated.xaml
loom generate:xaml contentview.swift --replace-region mainwindow.xaml --region-id main --overwrite
loom generate:swiftui mainwindow.xaml --view-name MainWindowScaffold --output MainWindowScaffold.swift
loom generate:juce contentview.swift --class-name MainComponent --output MainComponent.hpp
loom generate:contracts contentview.swift --target winui3 --json
loom patterns:transfer mainwindow.xaml --from winui3 --to macos --format json
loom patterns:transfer contentview.swift --from swiftui --to windows --format json
loom patterns:transfer mainwindow.qml --from qt --to windows --format json
loom patterns:transfer contentview.swift --from swiftui --to juce --format json
loom accessibility:audit mainwindow.xaml --fail-on warning
loom inspect:ascii mainwindow.xaml --output layout.txt --line-ending crlf
loom project:build examples/sampleapp/loom.json --output-dir examples/sampleapp/generated/project-build --overwrite --json
./examples/sampleapp/analyze-sample-app.sh --overwrite
```

## visual parity profiles

Use `inspect:font` to extract intrinsic font material properties from a supplied
TrueType, OpenType, TrueType Collection, or WOFF font file, or from an installed
family name. The JSON output includes parsed family names, ascender/descender
metrics, line gaps, cap height, x-height, weight/width classes,
italic/fixed-pitch data, normalized ratios, kerning table presence, and a
`profileTypography` block that can be copied into a visual parity profile.

`inspect:visual-parity` accepts a JSON profile that normalizes platform visual
defaults before comparison. It can also override either side's typography from
real font material with `--source-font`, `--target-font`, `--source-font-family`,
or `--target-font-family`.

Visual parity JSON includes per-node semantic descriptors, per-node
`provenance` for metrics, and per-finding confidence. Semantic descriptors flag
source-grounded differences in visual role, composition, chrome, emphasis,
density, interaction, media treatment, alignment, state, and native component
boundaries. This lets loom report cases where shared numeric constants match but
one platform still renders an old shell approximation, toolbar, media panel, tab
region, sidebar, or collection differently. Provenance marks values as `source`,
`font-material`,
`resolved-resource`, `style-setter`, `explicit-style-setter`, `profile`,
`resource-reference`, `default-profile`, or `unknown`, so reports distinguish
measured, explicitly provided, locally resolved, styled, referenced, and assumed
material properties. For XAML, loom resolves document-local resources plus local
merged dictionaries referenced with relative, absolute, or `ms-appx:///` paths,
implicit styles, explicit `Style="{StaticResource ...}"` references, and simple
`BasedOn` style chains. It also extracts common object-valued visual resources
and setters, including `Thickness`, `FontFamily`, and brush color values, and
follows simple `StaticResource`/`ThemeResource` alias chains. Profile files are
parsed strictly: unknown fields, unsupported platforms, and negative metrics are
rejected.

```json
{
  "schema_version": "1",
  "platforms": {
    "swiftui": {
      "typography": {
        "fontFamily": "Inter",
        "fallbackFonts": ["Segoe UI", "Arial"],
        "fontSize": 14,
        "kerning": 0,
        "lineHeight": 20,
        "baselineOffset": 0
      },
      "spacing": {
        "defaultPadding": 0,
        "defaultMargin": 0,
        "stackSpacing": 8
      },
      "controls": {
        "buttonMinHeight": 32,
        "textFieldMinHeight": 32,
        "toggleMinHeight": 32,
        "listRowMinHeight": 40
      }
    },
    "winui3": {
      "typography": {
        "fontFamily": "Inter",
        "fallbackFonts": ["Segoe UI", "Arial"],
        "fontSize": 14,
        "kerning": 0,
        "lineHeight": 20,
        "baselineOffset": 0
      },
      "spacing": {
        "defaultPadding": 0,
        "defaultMargin": 0,
        "stackSpacing": 8
      },
      "controls": {
        "buttonMinHeight": 32,
        "textFieldMinHeight": 32,
        "toggleMinHeight": 32,
        "listRowMinHeight": 40
      }
    }
  },
  "tolerances": {
    "distance": 0.5,
    "typography": 0.25
  }
}
```

## JSON/GCF interchange (1.1.0)

`capabilities [--format json|gcf|auto]` discovers formats and limits.
`data encode|decode|stats FILE|- [--format json|gcf|auto]` converts or measures
complete snapshots without executing them. Encode defaults to GCF; other data
commands default to JSON. JSON report commands accept GCF/Auto as alternatives;
report exit status, selected output paths and overwrite guards remain unchanged.
Manifests/profiles accept generic GCF with existing validation. Catalog files and
project-build artifacts retain their existing JSON/native source representation.
