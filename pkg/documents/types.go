package documents

import (
	"github.com/cheetahbyte/verba/pkg/bibmodel"
)

type DocumentMargin struct {
	Left   float64
	Right  float64
	Bottom float64
	Top    float64
}

type DocumentBibliography struct {
	Bibliography map[string]bibmodel.Citable // key → Eintrag
	Used         []bibmodel.Citable          // in Reihenfolge der Zitate
	CitationIDs  map[string]int              // key → ID [1], [2], ...
}

func (db *DocumentBibliography) GetEntry(key string) (bibmodel.Citable, bool) {
	entry, exists := db.Bibliography[key]
	return entry, exists
}

type Document struct {
	Margin       DocumentMargin
	PageWidth    float64
	Bibliography DocumentBibliography
}

func NewDocument() *Document {
	return &Document{
		Margin:    DocumentMargin{Left: 20, Right: 20, Top: 20, Bottom: 20},
		PageWidth: 210,
		Bibliography: DocumentBibliography{
			Bibliography: make(map[string]bibmodel.Citable),
			Used:         []bibmodel.Citable{},
			CitationIDs:  make(map[string]int),
		},
	}
}

func (d *Document) TextWidth() float64 {
	return d.PageWidth - d.Margin.Left - d.Margin.Right
}

func (d *Document) AddBibEntry(entry bibmodel.Citable) {
	key := entry.GetKey()

	if d.Bibliography.Bibliography == nil {
		d.Bibliography.Bibliography = make(map[string]bibmodel.Citable)
	}

	if _, exists := d.Bibliography.Bibliography[key]; !exists {
		d.Bibliography.Bibliography[key] = entry
	}
}

func (d *Document) Cite(key string) {
	if d.Bibliography.CitationIDs == nil {
		d.Bibliography.CitationIDs = make(map[string]int)
	}

	entry, exists := d.Bibliography.Bibliography[key]
	if !exists {
		return
	}

	if !d.Bibliography.IsUsed(entry) {
		d.Bibliography.Use(entry)

		newID := len(d.Bibliography.CitationIDs) + 1
		d.Bibliography.CitationIDs[key] = newID
	}
}

func (bib *DocumentBibliography) Use(entry bibmodel.Citable) {
	bib.Used = append(bib.Used, entry)
}

func (bib *DocumentBibliography) IsUsed(entry bibmodel.Citable) bool {
	key := entry.GetKey()
	for _, used := range bib.Used {
		if used.GetKey() == key {
			return true
		}
	}
	return false
}
