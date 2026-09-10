package loom

import (
	"fmt"
	"path/filepath"
	"strings"
)

func AnalyzeJUCE(path string) (Analysis, error) {
	data, err := readSourceFile(path, "JUCE source")
	if err != nil {
		return Analysis{}, fmt.Errorf("could not read JUCE source at %s: %w", path, err)
	}
	tokens, tokenDiagnostics := tokenizeSwift(string(data))
	diagnostics := append(juceTokenDiagnostics(tokenDiagnostics), juceDelimiterDiagnostics(tokens)...)
	children := parseJUCENodes(tokens)
	return Analysis{
		SourcePath:      path,
		RootView:        "juce",
		Component:       componentNameFromJUCESource(tokens, path),
		SyntaxNodeCount: countJUCEConstructTokens(tokens),
		Layout:          Node{Kind: KindRoot, Expression: "juce", Properties: map[string]string{"sourceDialect": "juce"}, Children: children},
		Diagnostics:     nonNilDiagnostics(diagnostics),
	}, nil
}

func parseJUCENodes(tokens []swiftToken) []Node {
	nodesByName := map[string]*Node{}
	order := []string{}
	for i := 0; i < len(tokens); i++ {
		name := juceConstructNameAt(tokens, i)
		if name == "" {
			continue
		}
		args := ""
		if j := nextNonNamespaceToken(tokens, i+1); j < len(tokens) && tokens[j].value == "(" {
			parser := &swiftParser{tokens: tokens, index: j}
			args = parser.consumeBalanced("(", ")")
			i = parser.index - 1
		}
		node := makeJUCENode(name, name, args)
		variable := juceAssignedIdentifier(tokens, i)
		if variable != "" {
			nodesByName[variable] = &node
			order = append(order, variable)
			continue
		}
		key := fmt.Sprintf("@%d", len(order))
		order = append(order, key)
		nodesByName[key] = &node
	}
	attached := map[string]bool{}
	for i := 0; i+2 < len(tokens); i++ {
		if tokens[i].value != "addAndMakeVisible" && tokens[i].value != "addChildComponent" {
			continue
		}
		j := nextNonNamespaceToken(tokens, i+1)
		if j >= len(tokens) || tokens[j].value != "(" {
			continue
		}
		parser := &swiftParser{tokens: tokens, index: j}
		args := parser.consumeBalanced("(", ")")
		name := firstJUCEIdentifierArgument(args)
		if node := nodesByName[name]; node != nil {
			attached[name] = true
		}
	}
	roots := []Node{}
	for _, key := range order {
		node := nodesByName[key]
		if node == nil {
			continue
		}
		if len(attached) > 0 && !attached[key] {
			continue
		}
		roots = append(roots, *node)
	}
	if len(roots) == 0 {
		for _, key := range order {
			if key == "" {
				continue
			}
			if node := nodesByName[key]; node != nil {
				roots = append(roots, *node)
			}
		}
	}
	return roots
}

func makeJUCENode(className, expression, args string) Node {
	props := map[string]string{"juceConstruct": className}
	node := Node{Kind: juceConstructKind(className), Expression: expression, Properties: props}
	if node.Kind == "" {
		node.Kind = KindComponent
		props["componentBoundary"] = "juce-component"
		props["requiresNativeImplementation"] = "true"
	}
	if args != "" {
		props["juce.arguments"] = args
		if text := firstSwiftStringArgument(args); text != "" {
			node.Arguments = quote(text)
			switch node.Kind {
			case KindText, KindButton, KindToggle:
				node.VisibleLabel = text
			case KindTextField:
				node.Placeholder = text
			case KindImage:
				node.Resource = text
				node.Arguments = quote("")
			}
		}
	}
	return node
}

func juceConstructKind(name string) NodeKind {
	name = strings.TrimPrefix(name, "juce::")
	switch name {
	case "Component", "Viewport":
		return KindVerticalStack
	case "FlexBox":
		return KindHorizontalStack
	case "Grid":
		return KindGrid
	case "TabbedComponent":
		return KindOverlayStack
	case "ListBox", "TableListBox", "TreeView":
		return KindList
	case "Label", "TextEditor":
		return KindText
	case "TextButton", "DrawableButton", "HyperlinkButton", "ToolbarButton":
		return KindButton
	case "Slider":
		return KindSlider
	case "ToggleButton":
		return KindToggle
	case "Image", "ImageComponent", "DrawableImage":
		return KindImage
	case "StretchableLayoutManager", "StretchableLayoutResizerBar":
		return KindSplitView
	default:
		return ""
	}
}

func juceConstructNameAt(tokens []swiftToken, index int) string {
	if index >= len(tokens) || tokens[index].kind != swiftTokenIdentifier {
		return ""
	}
	if index >= 2 && tokens[index-1].value == ":" && tokens[index-2].value == ":" {
		return ""
	}
	name := tokens[index].value
	if index+3 < len(tokens) && tokens[index+1].value == ":" && tokens[index+2].value == ":" && tokens[index+3].kind == swiftTokenIdentifier {
		name = tokens[index].value + "::" + tokens[index+3].value
	}
	if juceConstructKind(name) == "" {
		return ""
	}
	return name
}

func nextNonNamespaceToken(tokens []swiftToken, index int) int {
	if index+1 < len(tokens) && tokens[index].value == ":" && tokens[index+1].value == ":" {
		return index + 2
	}
	return index
}

func juceAssignedIdentifier(tokens []swiftToken, index int) string {
	for i := index - 1; i >= 0 && i >= index-6; i-- {
		if tokens[i].value == "=" || tokens[i].value == "new" {
			for j := i - 1; j >= 0 && j >= i-4; j-- {
				if tokens[j].kind == swiftTokenIdentifier && !isJUCETypeWord(tokens[j].value) {
					return tokens[j].value
				}
			}
		}
	}
	return ""
}

func isJUCETypeWord(value string) bool {
	switch value {
	case "auto", "new", "std", "unique_ptr", "make_unique", "juce":
		return true
	default:
		return juceConstructKind(value) != ""
	}
}

func firstJUCEIdentifierArgument(args string) string {
	tokens, _ := tokenizeSwift(args)
	for _, token := range tokens {
		if token.kind == swiftTokenIdentifier && token.value != "juce" {
			return token.value
		}
	}
	return ""
}

func countJUCEConstructTokens(tokens []swiftToken) int {
	count := 0
	for i := range tokens {
		if juceConstructNameAt(tokens, i) != "" {
			count++
		}
	}
	return count
}

func componentNameFromJUCESource(tokens []swiftToken, path string) string {
	for i := 0; i+2 < len(tokens); i++ {
		if tokens[i].value == "class" && tokens[i+1].kind == swiftTokenIdentifier && tokens[i+2].value == ":" {
			return tokens[i+1].value
		}
	}
	return componentName(path)
}

func looksLikeJUCESource(path string) bool {
	data, err := readSourceFile(path, "source")
	if err != nil {
		return false
	}
	source := string(data)
	return strings.Contains(source, "#include <JuceHeader.h>") ||
		strings.Contains(source, "#include \"JuceHeader.h\"") ||
		strings.Contains(source, "juce::") ||
		strings.Contains(strings.ToLower(filepath.Base(path)), "juce")
}

func juceDelimiterDiagnostics(tokens []swiftToken) []Diagnostic {
	for _, diagnostic := range swiftDelimiterDiagnostics(tokens) {
		diagnostic.Code = "JUCE.PARSE"
		diagnostic.Message = "unbalanced JUCE C++ source delimiters."
		return []Diagnostic{diagnostic}
	}
	return []Diagnostic{}
}

func juceTokenDiagnostics(in []Diagnostic) []Diagnostic {
	out := make([]Diagnostic, 0, len(in))
	for _, diagnostic := range in {
		diagnostic.Code = "JUCE.PARSE"
		diagnostic.Message = strings.ReplaceAll(diagnostic.Message, "Swift", "JUCE C++")
		out = append(out, diagnostic)
	}
	return out
}
