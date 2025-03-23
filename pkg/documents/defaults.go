package documents

var Article = DocumentMargin{
	Left:   25.0,
	Right:  25.0,
	Top:    30.0,
	Bottom: 30.0,
}

var Report = DocumentMargin{
	Left:   30.0,
	Right:  25.0,
	Top:    35.0,
	Bottom: 30.0,
}

var DocumentStyleClasses = map[string]DocumentMargin{
	"article": Article,
	"report":  Report,
}
