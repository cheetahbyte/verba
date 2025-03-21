package documents

var Article = Document{
	Margin: DocumentMargin{
		Left:   25.0,
		Right:  25.0,
		Top:    30.0,
		Bottom: 30.0,
	},
	PageWidth: 210.0,
}

var Report = Document{
	Margin: DocumentMargin{
		Left:   30.0,
		Right:  25.0,
		Top:    35.0,
		Bottom: 30.0,
	},
	PageWidth: 210.0,
}

var DocumentClasses = map[string]Document{
	"article": Article,
	"report":  Report,
}
