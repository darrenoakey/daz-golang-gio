package persist

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSaveAndLoadState(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("HOME", tmp)

	want := State{X: 100.5, Y: 200.5, Width: 1024, Height: 768}
	if err := SaveState("testapp", want); err != nil {
		t.Fatalf("SaveState: %v", err)
	}

	got, err := LoadState("testapp")
	if err != nil {
		t.Fatalf("LoadState: %v", err)
	}
	if got != want {
		t.Errorf("got %+v, want %+v", got, want)
	}
}

func TestLoadStateMissing(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("HOME", tmp)

	got, err := LoadState("nonexistent")
	if err != nil {
		t.Fatalf("LoadState should not error on missing file: %v", err)
	}
	if got != (State{}) {
		t.Errorf("expected zero state, got %+v", got)
	}
}

func TestLoadStateCorrupt(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("HOME", tmp)

	dir := filepath.Join(tmp, ".config", "daz-golang-gio")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "bad.json"), []byte("not json"), 0o644); err != nil {
		t.Fatal(err)
	}

	_, err := LoadState("bad")
	if err == nil {
		t.Fatal("expected error for corrupt JSON, got nil")
	}
}

func TestStatePathFormat(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("HOME", tmp)

	path := StatePath("myapp")
	want := filepath.Join(tmp, ".config", "daz-golang-gio", "myapp.json")
	if path != want {
		t.Errorf("StatePath = %q, want %q", path, want)
	}
}

func TestMultipleAppsIndependent(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("HOME", tmp)

	s1 := State{Width: 800, Height: 600, X: 10, Y: 20}
	s2 := State{Width: 1200, Height: 900, X: 30, Y: 40}

	if err := SaveState("app1", s1); err != nil {
		t.Fatal(err)
	}
	if err := SaveState("app2", s2); err != nil {
		t.Fatal(err)
	}

	got1, _ := LoadState("app1")
	got2, _ := LoadState("app2")

	if !got1.Equal(s1) {
		t.Errorf("app1: got %+v, want %+v", got1, s1)
	}
	if !got2.Equal(s2) {
		t.Errorf("app2: got %+v, want %+v", got2, s2)
	}
}

func TestSaveOverwrite(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("HOME", tmp)

	first := State{X: 10, Y: 20, Width: 300, Height: 400}
	if err := SaveState("overwrite", first); err != nil {
		t.Fatal(err)
	}

	second := State{X: 50, Y: 60, Width: 500, Height: 700}
	if err := SaveState("overwrite", second); err != nil {
		t.Fatal(err)
	}

	got, err := LoadState("overwrite")
	if err != nil {
		t.Fatal(err)
	}
	if !got.Equal(second) {
		t.Errorf("got %+v, want %+v", got, second)
	}
}

func TestStateValid(t *testing.T) {
	tests := []struct {
		name string
		s    State
		want bool
	}{
		{"zero", State{}, false},
		{"positive", State{Width: 100, Height: 100}, true},
		{"zero width", State{Width: 0, Height: 100}, false},
		{"negative height", State{Width: 100, Height: -1}, false},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.s.Valid(); got != tc.want {
				t.Errorf("Valid() = %v, want %v", got, tc.want)
			}
		})
	}
}
