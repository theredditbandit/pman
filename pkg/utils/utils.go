package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/charmbracelet/glamour"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/theredditbandit/pman/pkg"
	"github.com/theredditbandit/pman/pkg/db"
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
	err = filepath.Walk(pPath, func(_ string, info os.FileInfo, err error) error {
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
func BeautifyMD(data []byte) (string, error) {
	r, err := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(120),
		glamour.WithAutoStyle(),
	)
	if err != nil {
		return "", fmt.Errorf("something went wrong while creating renderer: %w", err)
	}
	out, err := r.Render(string(data))
	if err != nil {
		return "", err
	}
	return out, nil
}

// ReadREADME: returns the byte array of README.md of a project
func ReadREADME(projectName string) ([]byte, error) {
	actualName, err := db.GetRecord(projectName, pkg.ProjectAliasBucket)
	if err == nil {
		projectName = actualName
	}
	path, err := db.GetRecord(projectName, pkg.ProjectPaths)
	if err != nil {
		return nil, fmt.Errorf("project: %v not a valid project", projectName)
	}
	pPath := filepath.Join(path, "README.md")
	data, err := os.ReadFile(pPath)
	if err != nil {
		return nil, fmt.Errorf("Something went wrong while reading README for %s: %w", projectName, err)
	}
	return data, nil
}
