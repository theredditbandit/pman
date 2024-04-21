package pkg

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/theredditbandit/pman/pkg/db"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
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

func GetLastModifiedTime(pname string) string {
	var lastModTime time.Time
	var lastModFile string
	today := time.Now()
	_ = lastModFile
	pPath, err := db.GetRecord(pname, ProjectPaths)
	if err != nil {
		return "Something went wrong"
	}
	err = filepath.Walk(pPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && info.ModTime().After(lastModTime) {
			lastModTime = info.ModTime()
			lastModFile = info.Name()
		}
		return nil
	})
	if err != nil {
		return "Something went wrong"
	}
	switch fmt.Sprint(lastModTime.Date()) {
	case fmt.Sprint(today.Date()):
		return "Today"
	case fmt.Sprint(today.AddDate(0, 0, -1).Date()):
		return "Yesterday"
	}
	return fmt.Sprint(lastModTime.Date())
}
