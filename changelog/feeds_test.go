package changelog

import (
	"io/ioutil"
	"os"
	"testing"
)

var changelogs = []struct {
	File string
	Url  string
}{
	{
		File: "testdata/slackwareaarch64-current/ChangeLog.txt",
		Url:  "http://ftp.arm.slackware.com/slackwarearm/slackwareaarch64-current/ChangeLog.txt",
	},
	{
		File: "testdata/slackware64/ChangeLog.txt",
		Url:  "http://slackware.osuosl.org/slackware64-current/ChangeLog.txt",
	},
}

func TestFeed(t *testing.T) {
	for _, cl := range changelogs {
		fh, err := os.Open(cl.File)
		if err != nil {
			t.Fatal(err)
		}
		defer fh.Close()

		e, err := Parse(fh)
		if err != nil {
			t.Fatal(err)
		}

		f, err := ToFeed(cl.Url, e)
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
}
