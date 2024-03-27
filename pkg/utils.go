package pkg

import "strings"

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
	return strings.Title(s)
}
