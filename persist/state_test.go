package persist

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSaveAndLoadState(t *testing.T) {
	tmp := t.TempDir()
	origHome := os.Getenv("HOME")
	t.Setenv("HOME", tmp)
	t.Cleanup(func() { os.Setenv("HOME", origHome) })

	want := State{Width: 1024, Height: 768, X: 100, Y: 200, Mode: "windowed"}
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

	s1 := State{Width: 800, Height: 600}
	s2 := State{Width: 1200, Height: 900}

	if err := SaveState("app1", s1); err != nil {
		t.Fatal(err)
	}
	if err := SaveState("app2", s2); err != nil {
		t.Fatal(err)
	}

	got1, _ := LoadState("app1")
	got2, _ := LoadState("app2")

	if got1.Width != 800 || got2.Width != 1200 {
		t.Errorf("apps not independent: app1=%+v app2=%+v", got1, got2)
	}
}
