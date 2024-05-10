package db_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/theredditbandit/pman/pkg/db"
	bolt "go.etcd.io/bbolt"
)

func Test_GetDBLoc(t *testing.T) {
	t.Run("Test getDBLoc", func(t *testing.T) {
		expectedWords := []string{".local", "share", "pman"}

		actualPath, err := db.GetDBLoc(db.DBTestName)
		defer os.Remove(actualPath)

		require.Equal(t, err, nil)
		require.Contains(t, actualPath, expectedWords[0], expectedWords[1], expectedWords[2], db.DBTestName)
	})

	t.Run("Test GetDBLoc with empty dbname", func(t *testing.T) {
		dbname := ""
		expectedErr := db.ErrDBNameEmpty

		actualPath, err := db.GetDBLoc(dbname)

		require.Error(t, err)
		require.Equal(t, expectedErr, err)
		require.Empty(t, actualPath)
	})

	t.Run("Test getDBLoc with panic", func(t *testing.T) {

	})
}

func Test_GetRecord(t *testing.T) {
	t.Run("Test GetRecord", func(t *testing.T) {
		dbname := db.DBTestName
		DBLoc, err := db.GetDBLoc(dbname)
		require.NoError(t, err)
		defer os.Remove(DBLoc)

		expectedValue := "testValue"
		key := "testKey"
		bucketName := "testBucket"

		err = db.WriteToDB(dbname, map[string]string{key: expectedValue}, bucketName)
		require.NoError(t, err)

		actualValue, err := db.GetRecord(dbname, key, bucketName)
		require.NoError(t, err)
		require.Equal(t, expectedValue, actualValue)
	})

	t.Run("Test GetRecord with key not found", func(t *testing.T) {
		dbname := db.DBTestName
		DBLoc, err := db.GetDBLoc(dbname)
		require.NoError(t, err)
		defer os.Remove(DBLoc)
		key := "testKey"
		bucketName := "testBucket"
		expectedErr := db.ErrKeyNotFound

		err = db.WriteToDB(dbname, map[string]string{}, bucketName)
		require.NoError(t, err)

		actualValue, err := db.GetRecord(dbname, key, bucketName)

		require.Error(t, err)
		require.Equal(t, expectedErr, err)
		require.Empty(t, actualValue)
	})

	t.Run("Test GetRecord with bucket not found", func(t *testing.T) {
		dbname := db.DBTestName
		key := "testKey"
		bucketName := "testBucket"
		expectedErr := db.ErrBucketNotFound

		actualValue, err := db.GetRecord(dbname, key, bucketName)

		require.Error(t, err)
		require.Equal(t, expectedErr, err)
		require.Empty(t, actualValue)
	})
}
func Test_WriteToDB(t *testing.T) {
	t.Run("Test WriteToDB", func(t *testing.T) {
		dbname := db.DBTestName
		DBLoc, err := db.GetDBLoc(dbname)
		require.NoError(t, err)
		defer os.Remove(DBLoc)

		data := map[string]string{
			"key1": "value1",
			"key2": "value2",
			"key3": "value3",
		}
		bucketName := "testBucket"

		err = db.WriteToDB(dbname, data, bucketName)
		require.NoError(t, err)

		// Verify that the data was written correctly
		db, err := bolt.Open(DBLoc, 0o600, nil)
		require.NoError(t, err)
		defer db.Close()

		err = db.View(func(tx *bolt.Tx) error {
			bucket := tx.Bucket([]byte(bucketName))
			require.NotNil(t, bucket)

			for k, v := range data {
				value := bucket.Get([]byte(k))
				require.Equal(t, []byte(v), value)
			}

			return nil
		})
		require.NoError(t, err)
	})

	t.Run("Test WriteToDB with empty bucketname", func(t *testing.T) {
		dbname := db.DBTestName
		DBLoc, err := db.GetDBLoc(dbname)
		require.NoError(t, err)
		defer os.Remove(DBLoc)

		data := map[string]string{
			"key1": "value1",
			"key2": "value2",
			"key3": "value3",
		}
		bucketName := ""

		err = db.WriteToDB(dbname, data, bucketName)

		require.Error(t, err)
		require.ErrorIs(t, err, db.ErrCreateBucket)
	})

	t.Run("Test WriteToDB with empty map key", func(t *testing.T) {
		dbname := db.DBTestName
		DBLoc, err := db.GetDBLoc(dbname)
		require.NoError(t, err)
		defer os.Remove(DBLoc)

		data := map[string]string{
			"": "value1",
		}
		bucketName := "testBucket"

		err = db.WriteToDB(dbname, data, bucketName)

		require.Error(t, err)
		require.ErrorIs(t, err, db.ErrWriteToDB)
	})

	t.Run("Test WriteToDB with empty dbname value", func(t *testing.T) {
		dbname := ""
		data := map[string]string{
			"key1": "value1",
			"key2": "value2",
			"key3": "value3",
		}
		bucketName := "testBucket"

		err := db.WriteToDB(dbname, data, bucketName)

		require.Error(t, err)
		require.ErrorIs(t, err, db.ErrOpenDB)
	})
}

func Test_DeleteFromDb(t *testing.T) {
	t.Run("Test DeleteFromDb", func(t *testing.T) {
		dbname := db.DBTestName
		DBLoc, err := db.GetDBLoc(dbname)
		require.NoError(t, err)
		defer os.Remove(DBLoc)

		data := map[string]string{
			"key1": "value1",
			"key2": "value2",
			"key3": "value3",
		}
		bucketName := "testBucket"
		key := "key1"

		err = db.WriteToDB(dbname, data, bucketName)
		require.NoError(t, err)

		err = db.DeleteFromDb(dbname, key, bucketName)
		require.NoError(t, err)

		// Verify that the key was deleted
		db, err := bolt.Open(DBLoc, 0o600, nil)
		require.NoError(t, err)
		defer db.Close()

		err = db.View(func(tx *bolt.Tx) error {
			bucket := tx.Bucket([]byte(bucketName))
			require.NotNil(t, bucket)

			value := bucket.Get([]byte(key))
			require.Nil(t, value)

			return nil
		})
		require.NoError(t, err)
	})

	t.Run("Test DeleteFromDb with key not found", func(t *testing.T) {
		dbname := db.DBTestName
		DBLoc, err := db.GetDBLoc(dbname)
		require.NoError(t, err)
		defer os.Remove(DBLoc)

		data := map[string]string{
			"key1": "value1",
			"key2": "value2",
			"key3": "value3",
		}
		bucketName := "testBucket"
		key := "key4"

		err = db.WriteToDB(dbname, data, bucketName)
		require.NoError(t, err)

		err = db.DeleteFromDb(dbname, key, bucketName)

		require.NoError(t, err)
	})

	t.Run("Test DeleteFromDb with bucket not found", func(t *testing.T) {
		dbname := db.DBTestName
		key := "key1"
		bucketName := "testBucket"
		expectedErr := db.ErrBucketNotFound

		err := db.DeleteFromDb(dbname, key, bucketName)

		require.Error(t, err)
		require.ErrorIs(t, err, expectedErr)
	})
}

func Test_ListAllRecords(t *testing.T) {
	t.Run("Test ListAllRecords", func(t *testing.T) {
		dbname := db.DBTestName
		DBLoc, err := db.GetDBLoc(dbname)
		require.NoError(t, err)
		defer os.Remove(DBLoc)

		data := map[string]string{
			"key1": "value1",
			"key2": "value2",
			"key3": "value3",
		}
		bucketName := "testBucket"

		err = db.WriteToDB(dbname, data, bucketName)
		require.NoError(t, err)

		records, err := db.GetAllRecords(dbname, bucketName)

		require.NoError(t, err)
		require.Equal(t, data, records)
	})

	t.Run("Test ListAllRecords with bucket not found", func(t *testing.T) {
		dbname := db.DBTestName
		bucketName := "testBucket"
		expectedErr := db.ErrBucketNotFound

		records, err := db.GetAllRecords(dbname, bucketName)

		require.Error(t, err)
		require.ErrorIs(t, err, expectedErr)
		require.Nil(t, records)
	})
}
