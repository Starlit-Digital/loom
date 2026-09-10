# Screenshot Targets

These are the quickest snapshots to show Loom value in README, docs, or bug
reports. They are ordered as a beginner-to-reviewer story.

## Terminal Output Targets

Run these locally for clean, deterministic captures:

1. Health check: `loom status --json`
2. CLI readiness: `loom verify --json`
3. Catalog checks: `loom checks:command-catalog --json`
4. Parser evidence (SwiftUI): `loom inspect:swiftui examples/sampleapp/contentview.swift --json`
5. Parser evidence (XAML): `loom inspect:xaml examples/sampleapp/mainwindow.xaml --json`
6. ASCII layout review: `loom inspect:ascii examples/sampleapp/mainwindow.qml --output /tmp/loom-layout.txt`
7. Accessibility output: `loom accessibility:audit examples/sampleapp/mainwindow.xaml --format json --fail-on warning`
8. Contract report: `loom generate:contracts examples/sampleapp/contentview.swift --target winui3 --json`
9. Generator output: `loom generate:xaml examples/sampleapp/contentview.swift --output /tmp/generated.xaml`
10. Translator output: `loom generate:swiftui examples/sampleapp/mainwindow.xaml --view-name MainWindowScaffold --output /tmp/MainWindowScaffold.swift`
11. Transfer planning: `loom patterns:transfer examples/sampleapp/contentview.swift --from swiftui --to windows --format json`
12. Project build bundle: `loom project:build examples/sampleapp/loom.json --output-dir /tmp/loom-project-build --overwrite --json`

## UI/Website Targets

1. Project landing: https://sltd.ca/loom/
2. Release evidence panel (`docs/release-evidence.md` rendered on website)
3. Command examples section in README/docs/commands

## Terminal-only screenshot helper

Use this helper to capture a fixed set of terminal outputs for demos and PRs.

```sh
./scripts/capture-loom-screenshots.sh /tmp/loom-screenshots
```

By default, it writes plain-text output captures under `/tmp/loom-screenshots`,
then attempts a homepage screenshot if a browser automation tool is available.

## Optional visual workflow

When you want a quick hero image:

- Capture one clean homepage screenshot from desktop and mobile viewport.
- Capture one CLI window shot showing a successful `loom project:build` and
  concise JSON tail.
- Pair the image with the matching text transcript from `/tmp/loom-screenshots`.
