package db

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"

	bolt "go.etcd.io/bbolt"
)

const DBName = "projects.db"

// WriteToDB writes the data to the specified bucket in the database
func WriteToDB(data map[string]string, bucketName string) error {
	dbLoc, err := getDBLoc(DBName)
	if err != nil {
		return err
	}
	db, err := bolt.Open(dbLoc, 0o600, nil) // create the database if it doesn't exist then open it
	if err != nil {
		return err
	}
	defer db.Close()
	err = db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(bucketName))
		if err != nil {
			return err
		}
		for k, v := range data {
			err = bucket.Put([]byte(k), []byte(v))
			if err != nil {
				return err
			}
		}
		return nil
	})
	return err
}

func DeleteFromDb(key string, bucketName string) error {
	dbLoc, err := getDBLoc(DBName)
	if err != nil {
		return err
	}
	db, err := bolt.Open(dbLoc, 0o600, nil) // create the database if it doesn't exist then open it
	if err != nil {
		return err
	}
	defer db.Close()
	err = db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName))
		if bucket == nil {
			return fmt.Errorf("bucket %s not found", bucketName)
		}
		err := bucket.Delete([]byte(key))
		if err != nil {
			return err
		}
		return nil
	})
	return err
}

// getDBLoc returns the path to the database file, creating the directory if it doesn't exist
func getDBLoc(dbname string) (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	dbPath := filepath.Join(usr.HomeDir, ".local", "share", "pman", dbname)
	if _, err := os.Stat(filepath.Dir(dbPath)); os.IsNotExist(err) {
		err = os.MkdirAll(filepath.Dir(dbPath), 0o755)
		if err != nil {
			return "", err
		}
	}
	return dbPath, nil
}

// GetRecord returns the value of the key from the specified bucket, and error if it does not exist
func GetRecord(key string, bucketName string) (string, error) {
	dbLoc, err := getDBLoc(DBName)
	if err != nil {
		return "", err
	}
	var rec string
	db, err := bolt.Open(dbLoc, 0o600, nil)
	if err != nil {
		return "", err
	}
	defer db.Close()
	err = db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName))
		if bucket == nil {
			return fmt.Errorf("Bucket not found")
		}
		v := bucket.Get([]byte(key))
		if v == nil {
			return fmt.Errorf("Key not found in db\n")
		}
		rec = string(v)
		return nil
	})
	if err != nil {
		return "", err
	}
	return rec, nil
}

// GetAllRecords returns all the records from the specified bucket as a dictionary
func GetAllRecords(bucketName string) (map[string]string, error) {
	dbLoc, err := getDBLoc(DBName)
	if err != nil {
		return nil, err
	}
	db, err := bolt.Open(dbLoc, 0o600, nil)
	if err != nil {
		return nil, err
	}
	defer db.Close()
	records := make(map[string]string)
	err = db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName))
		if bucket == nil {
			return fmt.Errorf("Database not found. \nThis could be because no project dir has been initialized yet.")
		}
		err := bucket.ForEach(func(k, v []byte) error {
			records[string(k)] = string(v)
			return nil
		})
		return err
	})
	if err != nil {
		return nil, err
	}
	return records, nil
}

// UpdateRec updates the value of the key in the specified bucket, usually used to update the status of a project
func UpdateRec(key, val, bucketName string) error {
	dbLoc, err := getDBLoc(DBName)
	if err != nil {
		return err
	}
	db, err := bolt.Open(dbLoc, 0o600, nil)
	if err != nil {
		return err
	}
	defer db.Close()
	err = db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName))
		if bucket == nil {
			return fmt.Errorf("Bucket not found")
		}
		v := bucket.Get([]byte(key))
		if v == nil {
			return fmt.Errorf("Project not found")
		}
		err := bucket.Put([]byte(key), []byte(val))
		return err
	})
	return err
}

func DeleteDb() error {
	dbLoc, err := getDBLoc(DBName)
	if err != nil {
		return err
	}
	return os.Remove(dbLoc)
}
