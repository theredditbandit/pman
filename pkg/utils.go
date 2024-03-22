package pkg

import (
	"os"
	"os/user"
	"path/filepath"
)

// GetDB returns the path to the database file , creating the directory if it doesn't exist
func GetDB(dbname string) string {
	usr, err := user.Current()
	if err != nil {
		panic(err)
	}
	dbPath := filepath.Join(usr.HomeDir, ".local", "share", "projman", dbname)
	if _, err := os.Stat(filepath.Dir(dbPath)); os.IsNotExist(err) {
		os.MkdirAll(filepath.Dir(dbPath), 0755)
	}
	return dbPath
}
