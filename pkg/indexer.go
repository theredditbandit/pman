package pkg

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	c "github.com/theredditbandit/pman/constants"
	"github.com/theredditbandit/pman/pkg/db"
	"github.com/theredditbandit/pman/pkg/utils"
)

var (
	ErrDirname    = errors.New("error providing a directory name")
	ErrDirInvalid = errors.New("error providing a valid directory name")
	ErrIsNotDir   = errors.New("error providing a file instead of a directory")
	ErrIndexDir   = errors.New("error indexing directory")
)

// InitDirs indexes a directory for project directories and writes the data to the DB
func InitDirs(args []string) error {
	// the file which identifies a project directory
	if len(args) != 1 {
		log.Print("Please provide a directory name")
		return ErrDirname
	}
	dirname := args[0]
	if stat, err := os.Stat(dirname); os.IsNotExist(err) {
		log.Printf("%s is not a directory \n", dirname)
		return ErrDirInvalid
	} else if !stat.IsDir() {
		log.Printf("%s is a file and not a directory \n", dirname)
		return ErrIsNotDir
	}
	projDirs, err := indexDir(dirname)
	if err != nil {
		log.Print(err)
		return ErrIndexDir
	}
	fmt.Printf("Indexed %d project directories . . .\n", len(projDirs))
	projectStatusMap := make(map[string]string)
	projectPathMap := make(map[string]string)
	projectLastModTimeMap := make(map[string]string)
	for k, v := range projDirs { // k : full project path, v : project status ,
		projectName := filepath.Base(k)
		projectStatusMap[projectName] = v // filepath.Base(k) : project name
		projectPathMap[projectName] = k
	}
	err = db.WriteToDB(db.DBName, projectStatusMap, c.StatusBucket)
	if err != nil {
		log.Print(err)
		return err
	}
	err = db.WriteToDB(db.DBName, projectPathMap, c.ProjectPaths)
	if err != nil {
		log.Print(err)
		return err
	}

	for k := range projDirs {
		projectName := filepath.Base(k)
		t := utils.GetLastModifiedTime(db.DBName, projectName)
		lastEdited, timestamp := utils.ParseTime(t)
		projectLastModTimeMap[projectName] = fmt.Sprintf("%s-%d", lastEdited, timestamp)
	}

	err = db.WriteToDB(db.DBName, projectLastModTimeMap, c.LastUpdatedBucket)
	if err != nil {
		log.Print(err)
		return err
	}

	lastEdit := make(map[string]string)
	lastEdit["lastWrite"] = fmt.Sprint(time.Now().Format("02 Jan 06 15:04"))
	err = db.WriteToDB(db.DBName, lastEdit, c.ConfigBucket)
	if err != nil {
		log.Print(err)
		return err
	}

	return nil
}

// indexDir indexes a directory for project directories
func indexDir(path string) (map[string]string, error) {
	identifier := "readme"
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
		if !info.IsDir() {
			fileName := strings.ToLower(info.Name())
			if strings.Contains(fileName, identifier) {
				pname := filepath.Dir(path)
				absPath, _ := filepath.Abs(pname)
				projDirs[absPath] = "indexed"
				return filepath.SkipDir
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}
	return projDirs, nil
}
