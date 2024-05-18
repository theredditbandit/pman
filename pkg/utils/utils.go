package utils

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/charmbracelet/glamour"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/theredditbandit/pman/pkg"
	"github.com/theredditbandit/pman/pkg/db"
)

var (
	ErrBeautifyMD = errors.New("error beautifying markdown")
	ErrReadREADME = errors.New("error reading README")
)

func TitleCase(s string) string {
	c := cases.Title(language.English)
	return c.String(s)
}

func FilterByStatuses(data map[string]string, status []string) map[string]string {
	filteredData := make(map[string]string)
	for k, v := range data {
		for _, s := range status {
			if v == s {
				filteredData[k] = v
			}
		}
	}
	return filteredData
}

// Deprecated: Use ui.RenderTable instead
func PrintData(data map[string]string) {
	for k, v := range data {
		alias, err := db.GetRecord(db.DBName, k, pkg.ProjectAliasBucket)
		if err == nil {
			log.Printf("%s : %s (%s) \n", TitleCase(v), k, alias)
		} else {
			log.Printf("%s : %s  \n", TitleCase(v), k)
		}
	}
}

func GetLastModifiedTime(pname string) string {
	var lastModTime time.Time
	var lastModFile string
	today := time.Now()
	_ = lastModFile
	pPath, err := db.GetRecord(db.DBName, pname, pkg.ProjectPaths)
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
		log.Print("something went wrong while creating renderer: ", err)
		return "", errors.Join(ErrBeautifyMD, err)
	}
	out, _ := r.Render(string(data))
	return out, nil
}

// ReadREADME: returns the byte array of README.md of a project
func ReadREADME(projectName string) ([]byte, error) {
	path, err := db.GetRecord(db.DBName, projectName, pkg.ProjectPaths)
	if err != nil {
		actualName, err := db.GetRecord(db.DBName, projectName, pkg.ProjectAliasBucket)
		if err != nil {
			log.Printf("project: %v not a valid project\n", projectName)
			return nil, errors.Join(ErrReadREADME, err)
		}
		projectName = actualName
		path, err = db.GetRecord(db.DBName, projectName, pkg.ProjectPaths)
		if err != nil {
			log.Printf("project: %v not a valid project\n", projectName)
			return nil, errors.Join(ErrReadREADME, err)
		}
	}
	pPath := filepath.Join(path, "README.md")
	data, err := os.ReadFile(pPath)
	if err != nil {
		return nil, fmt.Errorf("something went wrong while reading README for %s: %w", projectName, err)
	}
	return data, nil
}
