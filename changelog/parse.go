package changelog

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strings"
	"time"
)

const (
	dividerStr     = `+--------------------------+`
	securityFixStr = `(* Security fix *)`
	dayPat         = `^(Mon|Tue|Wed|Thu|Fri|Sat|Sun)\s(Jan|Feb|Mar|Apr|May|Jun|Jul|Aug|Sep|Oct|Nov|Dec)\s.*\d{4}$`
	updatePat      = `^([a-z].*/.*):  (Added|Rebuilt|Removed|Updated|Upgraded)\.$`
)

var (
	dayReg    = regexp.MustCompile(dayPat)
	updateReg = regexp.MustCompile(updatePat)
)

// Parse takes in a slackware ChangeLog.txt and returns its collections of Entries
func Parse(r io.Reader) ([]Entry, error) {
	buf := bufio.NewReader(r)
	entries := []Entry{}
	curEntry := Entry{}
	var curUpdate *Update
	for {
		line, err := buf.ReadString('\n')
		if err != nil && err != io.EOF {
			return nil, err
		}
		isEOF := err == io.EOF
		trimmedline := strings.TrimSuffix(line, "\n")

		if trimmedline == dividerStr {
			if curUpdate != nil {
				curEntry.Updates = append(curEntry.Updates, *curUpdate)
				curUpdate = nil
			}
			entries = append(entries, curEntry)
			if isEOF {
				break
			}
			curEntry = Entry{}
		} else if dayReg.MatchString(trimmedline) {
			// this date means it is the beginning of an entry
			t, err := time.Parse(time.UnixDate, trimmedline)
			if err != nil {
				return nil, err
			}
			curEntry.Date = t
		} else if updateReg.MatchString(trimmedline) {
			// match on whether this is an update line
			if curUpdate != nil {
				curEntry.Updates = append(curEntry.Updates, *curUpdate)
				curUpdate = nil
			}
			m := updateReg.FindStringSubmatch(trimmedline)
			curUpdate = &Update{
				Name:   m[1],
				Action: m[2],
			}
		} else if curUpdate != nil && strings.HasPrefix(trimmedline, "  ") {
			curUpdate.Comment = curUpdate.Comment + line
		} else {
			// Everything else is a comment on the Entry
			curEntry.Comment = curEntry.Comment + line
		}

		if isEOF {
			break
		}
	}
	return entries, nil
}

// Entry is an section of updates (or release comments) in a ChangeLog.txt
type Entry struct {
	Date    time.Time
	Comment string
	Updates []Update
}

// SecurityFix is whether an update in this ChangeLog Entry includes a SecurityFix
func (e Entry) SecurityFix() bool {
	for _, u := range e.Updates {
		if u.SecurityFix() {
			return true
		}
	}
	return false
}

// ToChangeLog reformats the struct as the text for ChangeLog.txt output
func (e Entry) ToChangeLog() string {
	str := e.Date.Format(time.UnixDate) + "\n"
	if strings.Trim(e.Comment, " \n") != "" {
		str = str + e.Comment
	}
	for _, u := range e.Updates {
		str = str + u.ToChangeLog()
	}
	return str
}

// Update is a package or component that is updated in a ChangeLog Entry
type Update struct {
	Name    string
	Action  string
	Comment string
}

// SecurityFix that this update is a security fix (that the comment includes `(* Security fix *)`)
func (u Update) SecurityFix() bool {
	return strings.Contains(u.Comment, securityFixStr)
}

// ToChangeLog reformats the struct as the text for ChangeLog.txt output
func (u Update) ToChangeLog() string {
	return fmt.Sprintf("%s:  %s.\n%s", u.Name, u.Action, u.Comment)
}
