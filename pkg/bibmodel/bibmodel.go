package bibmodel

import "fmt"

type Citable interface {
	GetKey() string
	GetLabel() string        // z. B. für "[1]"
	FormatReference() string // formatiert für das Literaturverzeichnis
}

type BibEntry struct {
	Key    string
	Author string
	Title  string
	Year   string
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

type ArticleBibEntry struct {
	BibEntry
	Journal  string
	Number   string
	Pages    string
	DOI      string
	Keywords []string
}

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
