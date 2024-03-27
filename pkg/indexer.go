package pkg

import (
	"os"
	"path/filepath"
)
// IndexDir indexes a directory for multiple projects subdirectories
func IndexDir(path, identifier string) (map[string]int, error) { 
    // TODO : skip indexing the repo if it is already indexed in the db
	projDirs := make(map[string]int)
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() && info.Name() == ".git" {
			return filepath.SkipDir
		}
		if !info.IsDir() && info.Name() == identifier {
			pname := filepath.Dir(path)
			absPath, _ := filepath.Abs(pname)
			projDirs[absPath] = 0
			return filepath.SkipDir // skip the directory once we find the identifier
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return projDirs, nil
}