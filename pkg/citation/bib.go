package citation

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/cheetahbyte/verba/pkg/bibmodel"
	"github.com/cheetahbyte/verba/pkg/documents"
)

func LoadBibFile(path string, doc *documents.Document) error {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("Fehler beim Öffnen der Bib-Datei: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var entryType string
	var key string
	fields := map[string]string{}

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		if strings.HasPrefix(line, "@") {
			// Start eines neuen Eintrags
			entryType = ""
			key = ""
			fields = make(map[string]string)

			parts := strings.SplitN(line, "{", 2)
			if len(parts) == 2 {
				entryType = strings.ToLower(strings.TrimPrefix(parts[0], "@"))
				if comma := strings.Index(parts[1], ","); comma != -1 {
					key = strings.TrimSpace(parts[1][:comma])
				}
			}
		} else if strings.Contains(line, "=") {
			// Feld-Zuweisung (z. B. author = {...})
			parts := strings.SplitN(line, "=", 2)
			if len(parts) == 2 {
				field := strings.ToLower(strings.TrimSpace(parts[0]))
				value := strings.TrimSpace(parts[1])
				value = strings.Trim(value, "\",{}") // einfache Reinigung
				fields[field] = value
			}
		} else if line == "}" && key != "" {
			// Ende eines Eintrags → konkretes Struct erzeugen
			var entry bibmodel.Citable

			switch entryType {
			case "article":
				entry = &bibmodel.ArticleBibEntry{
					BibEntry: bibmodel.BibEntry{
						Key:    key,
						Author: fields["author"],
						Title:  fields["title"],
						Year:   fields["year"],
					},
					Journal:  fields["journal"],
					Number:   fields["number"],
					Pages:    fields["pages"],
					DOI:      fields["doi"],
					Keywords: strings.Split(fields["keywords"], ";"),
				}
			default:
				// Fallback: Basic BibEntry
				entry = &bibmodel.BibEntry{
					Key:    key,
					Author: fields["author"],
					Title:  fields["title"],
					Year:   fields["year"],
				}
			}

			doc.AddBibEntry(entry)
		}
	}

	return scanner.Err()
}
