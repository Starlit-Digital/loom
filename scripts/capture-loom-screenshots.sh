#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
OUT_DIR="${1:-/tmp/loom-screenshots}"
HOME_PAGE="${2:-https://sltd.ca/loom/}"

mkdir -p "${OUT_DIR}"

cd "${ROOT_DIR}"

loom_cmd() {
  go run ./cmd/loom "$@"
}

capture() {
  local label="$1"
  local command_text="$2"
  local outfile="${OUT_DIR}/${label}.txt"
  local exit_code

  echo "Capturing: ${command_text}" > "${outfile}"
  printf '%s\n' "" >> "${outfile}"
  if (cd "${ROOT_DIR}" && eval "${command_text}") >> "${outfile}" 2>&1; then
    exit_code=0
  else
    exit_code=$?
  fi
  echo "" >> "${outfile}"
  echo "exit_code=${exit_code}" >> "${outfile}"
}

capture "01-status" "loom_cmd status --json"
capture "02-verify" "loom_cmd verify --json"
capture "03-command-catalog" "loom_cmd checks:command-catalog --json"
capture "04-inspect-swiftui" "loom_cmd inspect:swiftui examples/sampleapp/contentview.swift --json"
capture "05-inspect-xaml" "loom_cmd inspect:xaml examples/sampleapp/mainwindow.xaml --json"
capture "06-inspect-ascii" "loom_cmd inspect:ascii examples/sampleapp/mainwindow.qml --output ${OUT_DIR}/layout-ascii.txt"
capture "07-a11y" "loom_cmd accessibility:audit examples/sampleapp/mainwindow.xaml --format json --fail-on warning"
capture "08-contracts" "loom_cmd generate:contracts examples/sampleapp/contentview.swift --target winui3 --json"
capture "09-generate-xaml" "loom_cmd generate:xaml examples/sampleapp/contentview.swift --output ${OUT_DIR}/generated-from-swiftui.xaml"
capture "10-generate-swiftui" "loom_cmd generate:swiftui examples/sampleapp/mainwindow.xaml --view-name MainWindowScaffold --output ${OUT_DIR}/generated-from-xaml.swift"
capture "11-transfer" "loom_cmd patterns:transfer examples/sampleapp/contentview.swift --from swiftui --to windows --format json"
capture "12-project-build" "loom_cmd project:build examples/sampleapp/loom.json --output-dir ${OUT_DIR}/project-build --overwrite --json"

if command -v npx >/dev/null 2>&1; then
  echo "Attempting homepage screenshot..."
  if npx playwright screenshot "${HOME_PAGE}" "${OUT_DIR}/loom-homepage.png" >/tmp/loom-playwright.log 2>&1; then
    echo "Homepage screenshot written to ${OUT_DIR}/loom-homepage.png"
  else
    echo "Homepage screenshot failed. Check /tmp/loom-playwright.log"
  fi
fi

echo "Done. Screenshot targets captured under ${OUT_DIR}."
