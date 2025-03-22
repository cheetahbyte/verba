package bibmodel

import "fmt"

type Citable interface {
	GetKey() string
	GetLabel() string
	FormatReference() string
	GetField(string) string
	GetType() string
	GetID() uint16
	SetID(uint16)
}

type BibEntry struct {
	Key    string
	Author string
	Title  string
	Year   string
	ID     uint16
}

func (a BibEntry) GetType() string {
	return "default"
}

func (a BibEntry) GetKey() string {
	return a.Key
}

func (a BibEntry) GetLabel() string {
	return "[" + a.Key + "]"
}

func (a BibEntry) FormatReference() string {
	return fmt.Sprintf("%s (%s). %s.",
		a.Author, a.Year, a.Title)
}

func (a BibEntry) GetID() uint16 {
	return a.ID
}

func (a *BibEntry) SetID(id uint16) {
	a.ID = id
}

func (b BibEntry) GetField(field string) string {
	switch field {
	case "key":
		return b.Key
	case "author":
		return b.Author
	case "title":
		return b.Title
	case "year":
		return b.Year
	}
	return ""
}

type ArticleBibEntry struct {
	BibEntry
	Journal  string
	Volume   string
	Number   string
	Pages    string
	DOI      string
	Keywords []string
}

func (a ArticleBibEntry) GetType() string { return "article" }

func (a ArticleBibEntry) GetKey() string {
	return a.Key
}

func (a ArticleBibEntry) GetLabel() string {
	return "[" + a.Key + "]"
}

func (a ArticleBibEntry) FormatReference() string {
	return fmt.Sprintf("%s (%s). %s. %s, S. %s. DOI: %s",
		a.Author, a.Year, a.Title, a.Journal, a.Pages, a.DOI)
}

func (a ArticleBibEntry) GetField(field string) string {
	switch field {
	case "journal":
		return a.Journal
	case "volume":
		return a.Volume
	case "number":
		return a.Number
	case "pages":
		return a.Pages
	case "doi":
		return a.DOI
	default:
		return a.BibEntry.GetField(field)
	}
}

func (a ArticleBibEntry) GetID() uint16 {
	return a.ID
}

func (a *ArticleBibEntry) SetID(id uint16) {
	a.ID = id
}
