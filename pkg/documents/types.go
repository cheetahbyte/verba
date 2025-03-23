package documents

type DocumentMargin struct {
	Left   float64
	Right  float64
	Bottom float64
	Top    float64
}
type Document struct {
	Margin    DocumentMargin
	PageWidth float64
}

func NewDocument() *Document {
	return &Document{
		Margin:    DocumentMargin{Left: 20, Right: 20, Top: 20, Bottom: 20},
		PageWidth: 210,
	}
}

func (d *Document) TextWidth() float64 {
	return d.PageWidth - d.Margin.Left - d.Margin.Right
}
