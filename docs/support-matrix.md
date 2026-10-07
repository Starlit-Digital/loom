# loom Support Matrix

loom `1.0.0` supports analyzer reports, transfer planning, component graphs,
reviewable generator scaffolds, target contracts, and manifest-directed project
build bundles across SwiftUI, WinUI XAML, Qt, and JUCE inputs.

## Current Analyzer Support

| Source dialect | Status | Current behavior |
| --- | --- | --- |
| WinUI XAML | Supported analyzer baseline | Parses common layout/control constructs, captures Grid track metadata, resolves selected visual resources and styles, audits transfer/accessibility risks. |
| SwiftUI | Supported analyzer baseline | Parses common layout/control constructs and visual modifiers into loom's shared layout model. |
| Qt QML | Supported analyzer baseline | Parses common Qt Quick and Controls layout/control constructs into loom's shared layout model. |
| Qt Designer UI | Conservative analyzer support | Parses common XML widget/layout constructs. |
| Qt C++ | Conservative analyzer support | Parses common layout/control construction heuristics. |
| JUCE C++ | Conservative analyzer and target support | Parses common `juce::Component` widget constructs, supports JUCE as a transfer/contract target, and emits reviewable JUCE Component scaffolds. |
| Fonts | Supported inspection | Extracts intrinsic OpenType/TrueType/TTC/WOFF metrics or installed family metrics for visual parity profiles. |

## Current JUCE Mapping

| JUCE construct | loom IR | Current behavior |
| --- | --- | --- |
| `juce::Component`, `juce::Viewport` | `verticalStack` / `scrollView` | Preserved as layout surfaces with explicit `resized()` policy in generated scaffolds. |
| `juce::FlexBox`, `juce::Grid` | `horizontalStack` / `grid` | Mapped as layout intent; exact imperative sizing remains a transfer policy. |
| `juce::Label` | `text` | Captures label-like text when available. |
| `juce::TextEditor` | `textField` | Captures editable text surface and reports state-binding contract needs. |
| `juce::TextButton`, `juce::DrawableButton`, `juce::HyperlinkButton` | `button` | Captures visible label when available and reports native action contract needs. |
| `juce::Slider` | `slider` | Captures range surface and reports state contract needs. |
| `juce::ToggleButton` | `toggle` | Captures binary state surface and reports state contract needs. |
| `juce::ListBox`, `juce::TableListBox`, `juce::TreeView` | `list` | Captures collection surface and reports collection/template contract needs. |
| `juce::ImageComponent`, `juce::DrawableImage` | `image` | Captures image surface intent. |
| Other JUCE components | `component` | Preserved as native JUCE component boundaries requiring review. |

## Current WinUI XAML Mapping

| WinUI XAML construct | loom IR | Current behavior |
| --- | --- | --- |
| `Grid` | `grid` | Parsed as a container. `Grid.RowDefinitions` and `Grid.ColumnDefinitions` are captured as metadata for transfer policy. |
| `StackPanel Orientation="Vertical"` | `verticalStack` | Parsed as a linear vertical layout. |
| `StackPanel Orientation="Horizontal"` | `horizontalStack` | Parsed as a linear horizontal layout. |
| `TextBlock` | `text` | Captures text content. |
| `Button`, `AppBarButton`, `HyperlinkButton` | `button` | Captures visible label when available and reports native action contract needs. |
| `TextBox`, `PasswordBox` | `textField` | Captures header or placeholder text and reports state-binding contract needs. |
| `Image` | `image` | Captures source and audits missing accessibility intent. |
| `ScrollViewer` | `scrollView` | Captures scroll container and audits nested scroll risk. |
| `ListView`, `GridView`, `ItemsRepeater` | `list` | Captures collection surface and reports collection/template contract needs. |
| `Slider`, `ProgressBar` | `slider` | Captures range/progress surface and reports state contract needs. |
| `ToggleSwitch`, `CheckBox` | `toggle` | Captures binary state surface and reports state contract needs. |
| `Rectangle Height="1"` / `Width="1"` | `divider` | Treated as a separator. |
| `Border Background=...` | `color` | Treated as a surface/color token candidate. |
| Other native WinUI controls | `component` | Preserved as unsupported native component boundaries with warnings and suggested fixes. |

## Generator And Translator Support

| Command | Current status | v1.0.0 expectation |
| --- | --- | --- |
| `patterns:transfer` | Implemented | Stable transfer planning across supported source/target pairs. |
| `inspect:parity` | Implemented | Stable structural parity report. |
| `inspect:visual-parity` | Implemented | Stable profile-normalized metric and semantic visual parity report. |
| `graph:components` | Implemented | Discover source-tree layout components and custom dependencies. |
| `generate:xaml` | Implemented scaffold | Emit reviewable WinUI XAML and support guarded owned-region replacement. |
| `generate:swiftui` | Implemented scaffold | Emit reviewable SwiftUI scaffolds. |
| `generate:juce` | Implemented scaffold | Emit reviewable JUCE Component C++ scaffolds. |
| `generate:contracts` | Implemented | Emit target native contracts for behavior, state, action, collection, component boundaries, and accessibility review. |
| `project:build` | Implemented bundle | Run manifest-directed validation, analysis, generated scaffold, contract, graph, transfer, parity, and summary workflows. |

Generated output is conservative scaffold output. Transfer and contract reports
identify structure, policy decisions, native behavior contracts, accessibility
gaps, and unsupported boundaries before a generated artifact is treated as
production code.
