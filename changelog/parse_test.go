package changelog

import (
	"os"
	"strings"
	"testing"
)

func TestParse(t *testing.T) {
	fh, err := os.Open("testdata/slackware64/ChangeLog.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer fh.Close()

	e, err := Parse(fh)
	if err != nil {
		t.Fatal(err)
	}

	// Make sure we got all the entries
	expectedLen := 52
	if len(e) != expectedLen {
		t.Errorf("expected %d entries; got %d", expectedLen, len(e))
	}

	// Make sure we got as many security fix entries as expected
	expectedSec := 34
	secCount := 0
	for i := range e {
		if e[i].SecurityFix() {
			secCount++
		}
	}
	if secCount != expectedSec {
		t.Errorf("expected %d security fix entries; got %d", expectedSec, secCount)
	}

	// Make sure we got as many individual updates as expected
	expectedUp := 597
	upCount := 0
	for i := range e {
		upCount += len(e[i].Updates)
	}
	if upCount != expectedUp {
		t.Errorf("expected %d updates across the entries; got %d", expectedUp, upCount)
	}

	// Make sure the top comment of an entry is working
	foundWorkmanComment := false
	expectedComment := "Thanks to Robby Workman for most of these updates."
	for i := range e {
		foundWorkmanComment = strings.Contains(e[i].Comment, expectedComment)
		if foundWorkmanComment {
			break
		}
	}
	if !foundWorkmanComment {
		t.Errorf("expected to find an Entry with comment %q", expectedComment)
	}
}
