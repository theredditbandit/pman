package db

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"

	bolt "go.etcd.io/bbolt"
)

const DBName string = "projects.db"
const DBTestName string = "projects_test.db"

var (
	ErrOpenDB          = errors.New("error opening database")
	ErrCreateBucket    = errors.New("error creating bucket")
	ErrWriteToDB       = errors.New("error writing to database")
	ErrBucketNotFound  = errors.New("bucket not found")
	ErrProjectNotFound = errors.New("project not found")
	ErrDeleteFromDB    = errors.New("error deleting from database")
	ErrKeyNotFound     = errors.New("key not found in db")
	ErrListAllRecords  = errors.New("error listing all records")
	ErrClearDB         = errors.New("error clearing database")
)

// WriteToDB writes the data to the specified bucket in the database
func WriteToDB(dbname string, data map[string]string, bucketName string) error {
	db, err := bolt.Open(getDBLoc(dbname), 0o600, nil) // create the database if it doesn't exist then open it
	if err != nil {
		log.Printf("%v : %v", ErrOpenDB, err)
		return errors.Join(ErrOpenDB, err)
	}
	defer db.Close()
	err = db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(bucketName))
		if err != nil {
			return errors.Join(ErrCreateBucket, err)
		}
		for k, v := range data {
			err = bucket.Put([]byte(k), []byte(v))
			if err != nil {
				return errors.Join(ErrWriteToDB, err)
			}
		}
		return nil
	})
	return err
}

func DeleteFromDb(dbname, key, bucketName string) error {
	db, err := bolt.Open(getDBLoc(dbname), 0o600, nil) // create the database if it doesn't exist then open it
	if err != nil {
		log.Printf("%v : %v", ErrOpenDB, err)
		return errors.Join(ErrOpenDB, err)
	}
	defer db.Close()
	err = db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName))
		if bucket == nil {
			return ErrBucketNotFound
		}
		err := bucket.Delete([]byte(key))
		if err != nil {
			return errors.Join(ErrDeleteFromDB, err)
		}
		return nil
	})
	return err
}

// getDBLoc returns the path to the database file, creating the directory if it doesn't exist
func getDBLoc(dbname string) string {
	usr, err := user.Current()
	if err != nil {
		panic(err)
	}
	dbPath := filepath.Join(usr.HomeDir, ".local", "share", "pman", dbname)
	if _, err := os.Stat(filepath.Dir(dbPath)); os.IsNotExist(err) {
		os.MkdirAll(filepath.Dir(dbPath), 0o755)
	}
	return dbPath
}

// GetRecord returns the value of the key from the specified bucket, and error if it does not exist
func GetRecord(dbname, key, bucketName string) (string, error) {
	var rec string
	db, err := bolt.Open(getDBLoc(dbname), 0o600, nil)
	if err != nil {
		log.Printf("%v : %v", ErrOpenDB, err)
		return "", errors.Join(ErrOpenDB, err)
	}
	defer db.Close()
	err = db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName))
		if bucket == nil {
			return ErrBucketNotFound
		}
		v := bucket.Get([]byte(key))
		if v == nil {
			return ErrKeyNotFound
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
func GetAllRecords(dbname, bucketName string) (map[string]string, error) {
	db, err := bolt.Open(getDBLoc(dbname), 0o600, nil)
	if err != nil {
		log.Printf("%v : %v", ErrOpenDB, err)
		return map[string]string{}, errors.Join(ErrOpenDB, err)
	}
	defer db.Close()
	records := make(map[string]string)
	err = db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName))
		if bucket == nil {
			fmt.Print("Database not found \nThis could be because no project dir has been initialized yet")
			return ErrBucketNotFound
		}
		err := bucket.ForEach(func(k, v []byte) error {
			records[string(k)] = string(v)
			return nil
		})
		return errors.Join(ErrListAllRecords, err)
	})
	if err != nil {
		return nil, err
	}
	return records, nil
}

// UpdateRec updates the value of the key in the specified bucket, usually used to update the status of a project
func UpdateRec(dbname, key, val, bucketName string) error {
	db, err := bolt.Open(getDBLoc(dbname), 0o600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	err = db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName))
		if bucket == nil {
			return ErrBucketNotFound
		}
		v := bucket.Get([]byte(key))
		if v == nil {
			return ErrProjectNotFound
		}
		err := bucket.Put([]byte(key), []byte(val))
		return err
	})
	return err
}

func DeleteDb(dbname string) error {
	err := os.Remove(getDBLoc(dbname))
	if err != nil {
		return errors.Join(ErrClearDB, err)
	}
	return nil
}
