// The pkg directory is being used to hold ancillary non-specific code, which could potentially be reused.
// A database model which could be used by other applications in the future (like a command line interface application) fits the bill here.
package models

import (
	"errors"
	"time"
)

var ErrNoRecord = errors.New("models: no matching record found")

// The fields pf the Snippet struct correspond to the fields in our MySQL snippets table
type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}
