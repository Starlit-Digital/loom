package loom

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestQMLLiteralsDoNotBecomeLayoutDelimiters(t *testing.T) {
	source := `import QtQuick
Item {
 property string closeBracket: "]"
 property string closeBrace: '}'
 function path(value) {
  const windows = /^\/[A-Za-z]:\//.test(value)
  const track = value.match(/^([0-9]+)\s*(?:[-_.)]\s*|\s+)/)
  const ratio = 12 / 3
  return windows ? '[' : value.replace(/[/*]/g, '')
 }
 Text { text: '(' }
 Badge {}
}`
	source = strings.Replace(source, "Badge {}", "property string templateValue: `hello ${name} { [] }`\n Badge {}", 1)
	tokens, diagnostics := tokenizeQML(source)
	if len(diagnostics) != 0 || len(qtDelimiterDiagnostics(tokens)) != 0 {
		t.Fatalf("valid QML diagnosed as unbalanced: %+v %+v", diagnostics, qtDelimiterDiagnostics(tokens))
	}
	refs := qmlComponentReferences(source)
	if len(refs) != 1 || refs[0] != "Badge" {
		t.Fatalf("custom references: %v", refs)
	}
}
func TestQMLMalformedDelimitersAndRegexRemainErrors(t *testing.T) {
	for _, source := range []string{`Item { Text { text: "ok" }`, "Item { property var pattern: /unterminated\n}", `Item { property string value: 'unterminated }`} {
		tokens, d := tokenizeQML(source)
		d = append(d, qtDelimiterDiagnostics(tokens)...)
		if len(d) == 0 {
			t.Fatalf("malformed QML accepted: %s", source)
		}
	}
}
func TestQMLGraphFindsCustomComponentBoundaries(t *testing.T) {
	root := t.TempDir()
	for name, source := range map[string]string{"Main.qml": "import QtQuick\nItem { property string ignored: 'Ghost {}'; Badge {} }", "Badge.qml": "import QtQuick\nRectangle {}"} {
		if err := os.WriteFile(filepath.Join(root, name), []byte(source), 0600); err != nil {
			t.Fatal(err)
		}
	}
	report, err := GraphComponents(root, "", "", "", "")
	if err != nil {
		t.Fatal(err)
	}
	found := false
	for _, edge := range report.Edges {
		if edge.From == "Main" && edge.To == "Badge" {
			found = true
		}
		if edge.To == "Ghost" {
			t.Fatal("string became dependency")
		}
	}
	if !found {
		t.Fatal("QML component reference missing", report.Edges)
	}
}
func TestSwiftQuotedDelimitersRemainLiterals(t *testing.T) {
	tokens, d := tokenizeSwift(`VStack { Text("]") }`)
	d = append(d, swiftDelimiterDiagnostics(tokens)...)
	if len(d) != 0 {
		t.Fatalf("quoted Swift delimiter: %+v", d)
	}
}
