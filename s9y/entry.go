package s9y

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

// returns the last_modified field (formatted) for the entry
func (e Entry) DateModified() string {
	if e.LastModified == 0 {
		return ""
	}
	return time.Unix(e.LastModified, 0).String()
}

// returns the timestamp field (publish date) for the entry
func (e Entry) DatePublished() string {
	return time.Unix(e.Timestamp, 0).String()
}

// re-create the old URL, used for files and alias
func (e Entry) URL() string {
	// pattern for URLs ID-Title.html
	var nonAlphanumericRegex = regexp.MustCompile(`[^a-zA-Z0-9\-\.\_ ]+`)

	var title = e.Title
	title = nonAlphanumericRegex.ReplaceAllString(title, "")
	title = strings.ReplaceAll(title, " ", "-")
	title = strings.ReplaceAll(title, "--", "-")

	return fmt.Sprintf("%d-%s", e.ID, title)
}
