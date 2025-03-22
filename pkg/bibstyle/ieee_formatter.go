package bibstyle

import (
	"fmt"

	"github.com/cheetahbyte/verba/pkg/bibmodel"
)

type IEEEFormatter struct{}

func (f IEEEFormatter) Format(entry bibmodel.Citable, id int) string {
	switch entry.GetType() {
	case "article":
		return fmt.Sprintf("[%d] %s. “%s.” %s, vol. %s, no. %s, pp. %s, %s. DOI: %s",
			id,
			entry.GetField("author"),
			entry.GetField("title"),
			entry.GetField("journal"),
			entry.GetField("volume"),
			entry.GetField("number"),
			entry.GetField("pages"),
			entry.GetField("year"),
			entry.GetField("doi"),
		)

	case "book":
		return fmt.Sprintf("[%d] %s. %s. %s, %s.",
			id,
			entry.GetField("author"),
			entry.GetField("title"),
			entry.GetField("publisher"),
			entry.GetField("year"),
		)

	default:
		return fmt.Sprintf("[%d] %s (%s). %s.",
			id,
			entry.GetField("author"),
			entry.GetField("year"),
			entry.GetField("title"),
		)
	}
}

func (f IEEEFormatter) CiteLabel(entry bibmodel.Citable) string {
	return fmt.Sprintf("[%d]", entry.GetID())
}
