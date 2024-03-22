package pkg

import (
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"

	bolt "go.etcd.io/bbolt"
)

func WriteToDB(projLocMap map[string]int, path string) error {
	db, err := bolt.Open(path, 0600, nil) // create the database if it doesn't exist then open it
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	err = db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte("projects"))
		if err != nil {
			return err
		}
		for k, v := range projLocMap {
			err = bucket.Put([]byte(k), []byte(fmt.Sprint(v)))
			if err != nil {
				return err
			}
		}
		return nil

	})
	return err
}

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


func GetRecord(key string) int {
	return 0
}

func Exists(key string) bool {
	return false
}
