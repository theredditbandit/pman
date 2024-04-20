package pkg

import (
	"fmt"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"pman/pkg/db"
)

func TitleCase(s string) string {
	c := cases.Title(language.English)
	return c.String(s)
}

func FilterByStatus(data map[string]string, status string) map[string]string {
	filteredData := make(map[string]string)
	for k, v := range data {
		if v == status {
			filteredData[k] = v
		}
	}
	return filteredData
}

// Deprecated: Use ui.RenderTable instead
func PrintData(data map[string]string) {
	for k, v := range data {
		alias, err := db.GetRecord(k, ProjectAliasBucket)
		if err == nil {
			fmt.Printf("%s : %s (%s) \n", TitleCase(v), k, alias)
		} else {
			fmt.Printf("%s : %s  \n", TitleCase(v), k)
		}
	}
}
