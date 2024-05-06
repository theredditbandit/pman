package pkg

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/theredditbandit/pman/pkg/db"
)

const (
	StatusBucket       = "projects"
	ProjectPaths       = "projectPaths"
	ProjectAliasBucket = "projectAliases"
)

// InitDirs indexes a directory for project directories and writes the data to the DB
func InitDirs(args []string) error {
	// the file which identifies a project directory
	projIdentifier := "README.md"
	if len(args) != 1 {
		return errors.New("Please provide a directory name")
	}
	dirname := args[0]
	if stat, err := os.Stat(dirname); os.IsNotExist(err) {
		return fmt.Errorf("%s is not a directory", dirname)
	} else if !stat.IsDir() {
		return fmt.Errorf("%s is a file and not a directory", dirname)
	}
	projDirs, err := indexDir(dirname, projIdentifier)
	if err != nil {
		return err
	}
	fmt.Printf("Indexed %d project directories . . .\n", len(projDirs))
	projectStatusMap := make(map[string]string)
	projectPathMap := make(map[string]string)
	for k, v := range projDirs { // k : full project path, v : project status ,
		projectStatusMap[filepath.Base(k)] = v // filepath.Base(k) : project name
		projectPathMap[filepath.Base(k)] = k
	}
	err = db.WriteToDB(projectStatusMap, StatusBucket)
	if err != nil {
		return err
	}
	err = db.WriteToDB(projectPathMap, ProjectPaths)
	if err != nil {
		return err
	}

	return nil
}

// indexDir indexes a directory for project directories
func indexDir(path, identifier string) (map[string]string, error) {
	projDirs := make(map[string]string)
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() && info.Name() == ".git" {
			pname := filepath.Dir(path)
			absPath, _ := filepath.Abs(pname)
			projDirs[absPath] = "indexed"
			return filepath.SkipDir
		}
		if !info.IsDir() && info.Name() == identifier {
			pname := filepath.Dir(path)
			absPath, _ := filepath.Abs(pname)
			projDirs[absPath] = "indexed"
			return filepath.SkipDir
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return projDirs, nil
}
