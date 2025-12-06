package testutil

import (
	"encoding/json"
	"flag"
	"os"
	"path/filepath"
	"testing"
)

var update = flag.Bool("update", false, "update golden files")

func AssertGoldenJSON(t *testing.T, golden string, got []byte) {
	t.Helper()

	if got == nil {
		return
	}

	// JSONの整形
	var gotJSON interface{}
	err := json.Unmarshal(got, &gotJSON)
	if err != nil {
		t.Fatalf("failed to unmarshal got JSON: %v", err)
	}

	got, err = json.MarshalIndent(gotJSON, "", "  ")

	if *update {
		err := os.MkdirAll(filepath.Dir(golden), 0755)
		err = os.WriteFile(golden, got, 0644)
		if err != nil {
			t.Fatalf("failed to write golden file: %v", err)
		}
	}

	want, err := os.ReadFile(golden)
	if err != nil {
		t.Fatalf("failed to read golden file: %v", err)
	}

	if string(got) != string(want) {
		t.Errorf("golden file mismatch:\nwant:\n%s\ngot:\n%s", string(want), string(got))
	}
}
