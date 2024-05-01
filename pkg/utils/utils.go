package utils

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/charmbracelet/glamour"
	"github.com/theredditbandit/pman/pkg"
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
		alias, err := db.GetRecord(k, pkg.ProjectAliasBucket)
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
	pPath, err := db.GetRecord(pname, pkg.ProjectPaths)
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
		return fmt.Sprintf("Today %s", lastModTime.Format("15:04"))
	case fmt.Sprint(today.AddDate(0, 0, -1).Date()):
		return fmt.Sprintf("Yesterday %s", lastModTime.Format("17:00"))
	}
	return fmt.Sprint(lastModTime.Format("02 Jan 06 15:04"))
}

// BeautifyMD: returns styled markdown
func BeautifyMD(data []byte) string {
	r, err := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(120),
		glamour.WithAutoStyle(),
	)
	if err != nil {
		log.Fatal("something went wrong while creating renderer: ", err)
	}
	out, _ := r.Render(string(data))
	return out
}

// ReadREADME: returns the byte array of REAMDE.md of a project
func ReadREADME(projectName string) []byte {
	actualName, err := db.GetRecord(projectName, pkg.ProjectAliasBucket)
	if err == nil {
		projectName = actualName
	}
	path, err := db.GetRecord(projectName, pkg.ProjectPaths)
	if err != nil {
		log.Fatalf("project: %v not a valid project\n", projectName)
	}
	pPath := filepath.Join(path, "README.md")
	data, err := os.ReadFile(pPath)
	if err != nil {
		log.Fatal("Something went wrong while reading README for ", projectName, "\nERR : ", err)
	}
	return data
}
