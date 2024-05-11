package pkg

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/theredditbandit/pman/pkg/db"
)

const (
	StatusBucket       string = "projects"
	ProjectPaths       string = "projectPaths"
	ProjectAliasBucket string = "projectAliases"
)

var (
	ErrDirname    error = errors.New("error providing a directory name")
	ErrDirInvalid error = errors.New("error providing a valid directory name")
	ErrIsNotDir   error = errors.New("error providing a file instead of a directory")
	ErrIndexDir   error = errors.New("error indexing directory")
)

// InitDirs indexes a directory for project directories and writes the data to the DB
func InitDirs(args []string) error {
	// the file which identifies a project directory
	projIdentifier := "README.md"
	if len(args) != 1 {
		fmt.Println("Please provide a directory name")
		return ErrDirname
	}
	dirname := args[0]
	if stat, err := os.Stat(dirname); os.IsNotExist(err) {
		fmt.Printf("%s is not a directory \n", dirname)
		return ErrDirInvalid
	} else if !stat.IsDir() {
		fmt.Printf("%s is a file and not a directory \n", dirname)
		return ErrIsNotDir
	}
	projDirs, err := indexDir(dirname, projIdentifier)
	if err != nil {
		fmt.Println(err)
		return ErrIndexDir
	}
	fmt.Printf("Indexed %d project directories . . .\n", len(projDirs))
	projectStatusMap := make(map[string]string)
	projectPathMap := make(map[string]string)
	for k, v := range projDirs { // k : full project path, v : project status ,
		projectStatusMap[filepath.Base(k)] = v // filepath.Base(k) : project name
		projectPathMap[filepath.Base(k)] = k
	}
	err = db.WriteToDB(db.DBName, projectStatusMap, StatusBucket)
	if err != nil {
		log.Print(err)
		return err
	}
	err = db.WriteToDB(db.DBName, projectPathMap, ProjectPaths)
	if err != nil {
		log.Print(err)
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
