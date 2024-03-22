package pkg

import (
	"fmt"
	"log"

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
