package fetch

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestFetchChangeLog(t *testing.T) {
	server := httptest.NewServer(http.FileServer(http.Dir("../changelog/testdata/slackware64/")))
	defer server.Close()

	r := Repo{
		URL: server.URL,
	}

	e, mtime, err := r.ChangeLog()
	if err != nil {
		t.Fatal(err)
	}

	expectedLen := 52
	if len(e) != expectedLen {
		t.Errorf("expected %d entries; got %d", expectedLen, len(e))
	}

	stat, err := os.Stat("../changelog/testdata/slackware64/ChangeLog.txt")
	if err != nil {
		t.Fatal(err)
	}

	if mtime.Unix() != stat.ModTime().Unix() {
		t.Errorf("time stamps not the same: expected %d; got %d", stat.ModTime().Unix(), mtime.Unix())
	}
}
