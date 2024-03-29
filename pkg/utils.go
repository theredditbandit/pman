package pkg

import (
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var SupportedStatus = []string{
	"Indexed",
	"Started",
	"Paused",
	"Completed",
	"Aborted",
	"Deleted",
	"Ongoing",
	"Not Started",
}

func Title(s string) string {
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

func PrintData(data map[string]string) {
	for k, v := range data {
		println(Title(v), " : ", k)
	}
}
