package changelog

import (
	"fmt"
	"time"

	"github.com/gorilla/feeds"
)

// ToFeed produces a github.com/gorilla/feeds.Feed that can be written to Atom or Rss
func ToFeed(link string, entries []Entry) (*feeds.Feed, error) {
	var newestEntryTime time.Time
	var oldestEntryTime time.Time

	for _, e := range entries {
		if e.Date.After(newestEntryTime) {
			newestEntryTime = e.Date
		}
		if e.Date.Before(oldestEntryTime) {
			oldestEntryTime = e.Date
		}
	}

	feed := &feeds.Feed{
		Title:       "",
		Link:        &feeds.Link{Href: link},
		Description: "Generated ChangeLog.txt feeds by sl-feeds (github.com/vbatts/sl-feeds)",
		Created:     oldestEntryTime,
		Updated:     newestEntryTime,
	}
	feed.Items = make([]*feeds.Item, len(entries))
	for i, e := range entries {
		feed.Items[i] = &feeds.Item{
			Created:     e.Date,
			Link:        &feeds.Link{Href: ""},
			Description: e.ToChangeLog(),
		}

		updateWord := "updates"
		if len(e.Updates) == 1 {
			updateWord = "update"
		}
		if e.SecurityFix() {
			feed.Items[i].Title = fmt.Sprintf("%d %s. Including a %s!", len(e.Updates), updateWord, securityFixStr)
		} else {
			feed.Items[i].Title = fmt.Sprintf("%d %s.", len(e.Updates), updateWord)
		}
	}

	return feed, nil
}