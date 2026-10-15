package loom

import (
	"strings"
	"unicode"
)

// tokenizeQML recognizes QML's embedded JavaScript literals. It is a bounded
// structural lexer, not a complete ECMAScript parser; template bodies stay opaque.
func tokenizeQML(source string) ([]swiftToken, []Diagnostic) {
	runes := []rune(source)
	tokens := []swiftToken{}
	diagnostics := []Diagnostic{}
	diagnostic := func(start int, what string) {
		offset := start
		diagnostics = append(diagnostics, Diagnostic{Severity: SeverityError, Code: "QT.PARSE", Message: "unterminated Qt " + what + ".", SourceOffset: &offset})
	}
	for i := 0; i < len(runes); {
		r := runes[i]
		if unicode.IsSpace(r) {
			i++
			continue
		}
		if r == '/' && i+1 < len(runes) && runes[i+1] == '/' {
			for i < len(runes) && runes[i] != '\n' {
				i++
			}
			continue
		}
		if r == '/' && i+1 < len(runes) && runes[i+1] == '*' {
			start := i
			i += 2
			closed := false
			for i+1 < len(runes) {
				if runes[i] == '*' && runes[i+1] == '/' {
					i += 2
					closed = true
					break
				}
				i++
			}
			if !closed {
				diagnostic(start, "block comment")
				i = len(runes)
			}
			continue
		}
		if r == '"' || r == '\'' || r == '`' {
			start := i
			quote := r
			i++
			var value strings.Builder
			closed := false
			for i < len(runes) {
				if runes[i] == '\\' && i+1 < len(runes) {
					value.WriteRune(runes[i+1])
					i += 2
					continue
				}
				if runes[i] == quote {
					i++
					closed = true
					break
				}
				value.WriteRune(runes[i])
				i++
			}
			if !closed {
				diagnostic(start, "string literal")
			}
			tokens = append(tokens, swiftToken{kind: swiftTokenString, value: value.String(), offset: start})
			continue
		}
		if r == '/' && qmlRegexCanStart(tokens) {
			start := i
			i++
			inClass := false
			closed := false
			for i < len(runes) && runes[i] != '\n' && runes[i] != '\r' {
				if runes[i] == '\\' && i+1 < len(runes) {
					i += 2
					continue
				}
				if runes[i] == '[' {
					inClass = true
				} else if runes[i] == ']' {
					inClass = false
				} else if runes[i] == '/' && !inClass {
					i++
					closed = true
					break
				}
				i++
			}
			if closed {
				for i < len(runes) && unicode.IsLetter(runes[i]) {
					i++
				}
			} else {
				diagnostic(start, "regular expression")
			}
			tokens = append(tokens, swiftToken{kind: swiftTokenString, value: string(runes[start:i]), offset: start})
			continue
		}
		if isSwiftIdentifierStart(r) {
			start := i
			i++
			for i < len(runes) && isSwiftIdentifierPart(runes[i]) {
				i++
			}
			tokens = append(tokens, swiftToken{kind: swiftTokenIdentifier, value: string(runes[start:i]), offset: start})
			continue
		}
		if unicode.IsDigit(r) {
			start := i
			i++
			for i < len(runes) && (unicode.IsDigit(runes[i]) || runes[i] == '.') {
				i++
			}
			tokens = append(tokens, swiftToken{kind: swiftTokenSymbol, value: string(runes[start:i]), offset: start})
			continue
		}
		tokens = append(tokens, swiftToken{kind: swiftTokenSymbol, value: string(r), offset: i})
		i++
	}
	return tokens, diagnostics
}
func qmlRegexCanStart(tokens []swiftToken) bool {
	if len(tokens) == 0 {
		return true
	}
	last := tokens[len(tokens)-1]
	if last.kind == swiftTokenIdentifier {
		switch last.value {
		case "return", "throw", "case", "delete", "void", "typeof", "instanceof", "in", "of", "yield", "await":
			return true
		}
		return false
	}
	if last.kind != swiftTokenSymbol {
		return false
	}
	return strings.Contains("=:(,[! ?{;+*- %&|^~<>", last.value) && len(last.value) == 1
}

func qmlComponentReferences(source string) []string {
	tokens, _ := tokenizeQML(source)
	seen := map[string]bool{}
	out := []string{}
	for i, t := range tokens {
		if t.kind != swiftTokenIdentifier || i+1 >= len(tokens) || tokens[i+1].kind != swiftTokenSymbol || tokens[i+1].value != "{" || t.value == "" || !unicode.IsUpper([]rune(t.value)[0]) {
			continue
		}
		if qtConstructKind(t.value) != "" || qmlFrameworkType(t.value) || seen[t.value] {
			continue
		}
		seen[t.value] = true
		out = append(out, t.value)
	}
	return out
}
func qmlFrameworkType(name string) bool {
	switch name {
	case "Item", "QtObject", "Component", "Connections", "Timer", "MouseArea", "TapHandler", "HoverHandler", "DragHandler", "WheelHandler", "Shortcut", "Window", "ApplicationWindow", "Popup", "Dialog", "DialogButtonBox", "Action", "ActionGroup", "Menu", "MenuBar", "MenuItem", "MenuSeparator", "ToolTip", "ButtonGroup", "Button", "Pane", "Frame", "GroupBox", "Control", "Container", "Page", "PageIndicator", "ToolBar", "TabBar", "TabButton", "ComboBox", "SpinBox", "Dial", "RangeSlider", "BusyIndicator", "ScrollBar", "ScrollIndicator", "TextMetrics", "FontMetrics", "FontLoader", "Binding", "PropertyChanges", "AnchorChanges", "ParentChange", "State", "Transition", "Behavior", "NumberAnimation", "ColorAnimation", "PropertyAnimation", "ParallelAnimation", "SequentialAnimation", "PauseAnimation", "ScriptAction", "ShaderEffect", "ShaderEffectSource", "Canvas", "Gradient", "GradientStop", "Loader", "MultiEffect", "DropShadow", "ColorOverlay", "FileDialog", "FolderDialog", "Settings", "Screen", "SystemPalette", "StandardPaths", "MediaPlayer", "AudioOutput", "VideoOutput", "SoundEffect", "ListModel", "ListElement":
		return true
	}
	return false
}
