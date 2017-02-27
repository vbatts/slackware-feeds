package changelog

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestFeed(t *testing.T) {
	fh, err := os.Open("testdata/slackware64/ChangeLog.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer fh.Close()

	e, err := Parse(fh)
	if err != nil {
		t.Fatal(err)
	}

	f, err := ToFeed("http://slackware.osuosl.org/slackware64-current/ChangeLog.txt", e)
	if err != nil {
		t.Fatal(err)
	}

	rss, err := f.ToRss()
	if err != nil {
		t.Fatal(err)
	}
	//println(rss)
	if len(rss) == 0 {
		t.Error("rss output is empty")
	}

	if err := f.WriteRss(ioutil.Discard); err != nil {
		t.Error(err)
	}
}
