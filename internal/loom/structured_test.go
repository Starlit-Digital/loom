package loom

import (
	"bytes"
	"encoding/json"
	"github.com/cshaiku/loom/internal/structured"
	"os"
	"path/filepath"
	"testing"
)

func TestGCFReportsAndOverwriteGuard(t *testing.T) {
	path := "../../examples/sampleapp/mainwindow.xaml"
	var legacy, out, errout bytes.Buffer
	if err := Run([]string{"inspect:xaml", path, "--format", "json"}, &legacy, &errout); err != nil {
		t.Fatal(err)
	}
	if err := Run([]string{"inspect:xaml", path, "--format", "gcf"}, &out, &errout); err != nil {
		t.Fatal(err)
	}
	decoded, err := structured.JSON(out.Bytes(), structured.MaxBytes)
	if err != nil {
		t.Fatal(err)
	}
	var a, b any
	json.Unmarshal(legacy.Bytes(), &a)
	json.Unmarshal(decoded, &b)
	aa, _ := json.Marshal(a)
	bb, _ := json.Marshal(b)
	if !bytes.Equal(aa, bb) {
		t.Fatal("changed report semantics")
	}
	output := filepath.Join(t.TempDir(), "report.gcf")
	args := []string{"inspect:xaml", path, "--format", "gcf", "--output", output}
	if err = Run(args, &out, &errout); err != nil {
		t.Fatal(err)
	}
	original, _ := os.ReadFile(output)
	if !bytes.HasPrefix(original, []byte("GCF")) {
		t.Fatal("file is not GCF")
	}
	if err = Run(args, &out, &errout); err == nil {
		t.Fatal("overwrite guard bypassed")
	}
	after, _ := os.ReadFile(output)
	if !bytes.Equal(original, after) {
		t.Fatal("guard changed file")
	}
}
func TestGCFManifestProfileAndDiscovery(t *testing.T) {
	for _, name := range []string{"loom.json", "visual-profile.json"} {
		original, err := os.ReadFile("../../examples/sampleapp/" + name)
		if err != nil {
			t.Fatal(err)
		}
		var value any
		json.Unmarshal(original, &value)
		wire, err := structured.Marshal(value, "gcf")
		if err != nil {
			t.Fatal(err)
		}
		path := filepath.Join(t.TempDir(), name+".gcf")
		os.WriteFile(path, wire, 0600)
		if name == "loom.json" {
			if _, err = readLoomManifest(path); err != nil {
				t.Fatal(err)
			}
		} else {
			if _, err = LoadVisualProfile(path); err != nil {
				t.Fatal(err)
			}
		}
	}
	var out, errout bytes.Buffer
	if err := Run([]string{"capabilities", "--format", "gcf"}, &out, &errout); err != nil {
		t.Fatal(err)
	}
	if _, err := structured.JSON(out.Bytes(), structured.MaxBytes); err != nil {
		t.Fatal(err)
	}
	out.Reset()
	if err := Run([]string{"list", "--format", "gcf"}, &out, &errout); err != nil {
		t.Fatal(err)
	}
	if !bytes.HasPrefix(out.Bytes(), []byte("GCF")) {
		t.Fatal("list not encoded")
	}
}
