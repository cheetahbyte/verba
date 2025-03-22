package bibstyle

import (
	"fmt"
	"strings"

	"github.com/cheetahbyte/verba/pkg/bibmodel"
)

type APAFormatter struct{}

func (f APAFormatter) Format(entry bibmodel.Citable, _ int) string {
	switch entry.GetType() {
	case "article":
		return fmt.Sprintf("%s (%s). %s. %s.",
			entry.GetField("author"),
			entry.GetField("year"),
			entry.GetField("title"),
			entry.GetField("journal"),
		)

	case "book":
		return fmt.Sprintf("%s (%s). %s. %s.",
			entry.GetField("author"),
			entry.GetField("year"),
			entry.GetField("title"),
			entry.GetField("publisher"),
		)

	default:
		return fmt.Sprintf("%s (%s). %s.",
			entry.GetField("author"),
			entry.GetField("year"),
			entry.GetField("title"),
		)
	}
}

func (f APAFormatter) CiteLabel(entry bibmodel.Citable) string {
	authors := getAuthors(entry.GetField("author"))
	if len(authors) == 1 {
		return fmt.Sprintf("(%s, %s)", entry.GetField("author"), entry.GetField("year"))
	}
	return fmt.Sprintf("(%s, %s)", strings.Split(authors[0], ",")[0]+" et.al.", entry.GetField("year"))
}
