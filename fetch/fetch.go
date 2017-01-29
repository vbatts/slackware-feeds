package fetch

import (
	"fmt"
	"net/http"
	"time"

	"../changelog"
)

// Repo represents a remote slackware software repo
type Repo struct {
	URL string
}

func (r Repo) get(file string) (*http.Response, error) {
	return http.Get(r.URL + "/" + file)
}

// ChangeLog fetches the ChangeLog.txt for this remote Repo, along with the
// last-modified (for comparisons).
func (r Repo) ChangeLog() (e []changelog.Entry, mtime time.Time, err error) {
	resp, err := r.get("ChangeLog.txt")
	if err != nil {
		return nil, time.Unix(0, 0), err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, time.Unix(0, 0), fmt.Errorf("%d status from %s", resp.StatusCode, resp.Request.URL)
	}

	mtime, err = http.ParseTime(resp.Header.Get("last-modified"))
	if err != nil {
		return nil, time.Unix(0, 0), err
	}

	e, err = changelog.Parse(resp.Body)
	if err != nil {
		return nil, mtime, err
	}
	return e, mtime, nil
}
