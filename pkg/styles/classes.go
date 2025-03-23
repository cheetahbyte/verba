package styles

import "github.com/cheetahbyte/verba/pkg/documents"

var Article = documents.DocumentMargin{
	Left:   25.0,
	Right:  25.0,
	Top:    30.0,
	Bottom: 30.0,
}

var Report = documents.DocumentMargin{
	Left:   30.0,
	Right:  25.0,
	Top:    35.0,
	Bottom: 30.0,
}

var DocumentStyleClasses = map[string]documents.DocumentMargin{
	"article": Article,
	"report":  Report,
}
