package util

import "testing"

func TestFind(t *testing.T) {
	paths, err := FindFiles("../changelog", "ChangeLog.txt")
	if err != nil {
		t.Fatal(err)
	}

	if len(paths) != 1 {
		t.Errorf("expected to find %d file, but found %d", 1, len(paths))
	}
}
