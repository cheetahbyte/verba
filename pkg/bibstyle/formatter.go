package bibstyle

import (
	"strings"

	"github.com/cheetahbyte/verba/pkg/bibmodel"
)

type Formatter interface {
	Format(entry bibmodel.Citable, id int) string
	CiteLabel(entry bibmodel.Citable) string
}

// used for et.al.
func extractfirstAuthors(authorsString string) string {
	authors := strings.Split(authorsString, "and")
	return strings.Trim(authors[0], " ")
}

func getAuthors(authorsString string) []string {
	return strings.Split(authorsString, "and")
}
