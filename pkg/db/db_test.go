package db_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/theredditbandit/pman/pkg/db"
	bolt "go.etcd.io/bbolt"
)

const dbname = db.DBTestName
const bucketName = "testBucket"
const key = "testKey"

func Test_GetDBLoc(t *testing.T) {
	t.Run("Test getDBLoc", func(t *testing.T) {
		expectedWords := []string{".local", "share", "pman"}

		actualPath, err := db.GetDBLoc(dbname)

		t.Cleanup(func() {
			os.Remove(actualPath)
		})

		assert.NoError(t, err)
		assert.Contains(t, actualPath, expectedWords[0])
		assert.Contains(t, actualPath, expectedWords[1])
		assert.Contains(t, actualPath, expectedWords[2])
		assert.Contains(t, actualPath, db.DBTestName)
	})

	t.Run("Test GetDBLoc with empty dbname", func(t *testing.T) {
		dbname := ""
		expectedErr := db.ErrDBNameEmpty

		actualPath, err := db.GetDBLoc(dbname)

		assert.ErrorIs(t, err, expectedErr)
		assert.Empty(t, actualPath)
	})
}

func Test_GetRecord(t *testing.T) {
	t.Run("Test GetRecord", func(t *testing.T) {
		actualPath, err := db.GetDBLoc(dbname)
		assert.NoError(t, err)

		t.Cleanup(func() {
			os.Remove(actualPath)
		})

		expectedValue := "testValue"

		err = db.WriteToDB(dbname, map[string]string{key: expectedValue}, bucketName)
		assert.NoError(t, err)

		actualValue, err := db.GetRecord(dbname, key, bucketName)
		assert.NoError(t, err)
		assert.Equal(t, expectedValue, actualValue)
	})

	t.Run("Test GetRecord with empty dbname", func(t *testing.T) {
		dbname := ""
		key := "testKey"
		bucketName := "testBucket"
		expectedErr := db.ErrDBNameEmpty

		actualValue, err := db.GetRecord(dbname, key, bucketName)

		require.ErrorIs(t, err, expectedErr)
		assert.Empty(t, actualValue)
	})

	t.Run("Test GetRecord with key not found", func(t *testing.T) {
		expectedErr := db.ErrKeyNotFound
		actualPath, err := db.GetDBLoc(dbname)
		assert.NoError(t, err)

		t.Cleanup(func() {
			os.Remove(actualPath)
		})

		err = db.WriteToDB(dbname, map[string]string{}, bucketName)
		assert.NoError(t, err)

		actualValue, err := db.GetRecord(dbname, key, bucketName)

		require.ErrorIs(t, err, expectedErr)
		assert.Empty(t, actualValue)
	})

	t.Run("Test GetRecord with bucket not found", func(t *testing.T) {
		expectedErr := db.ErrBucketNotFound

		actualValue, err := db.GetRecord(dbname, key, bucketName)

		require.ErrorIs(t, err, expectedErr)
		assert.Empty(t, actualValue)
	})
}
func Test_WriteToDB(t *testing.T) {
	t.Run("Test WriteToDB", func(t *testing.T) {
		actualPath, err := db.GetDBLoc(dbname)
		assert.NoError(t, err)

		t.Cleanup(func() {
			os.Remove(actualPath)
		})

		data := map[string]string{
			"key1": "value1",
			"key2": "value2",
			"key3": "value3",
		}
		bucketName := "testBucket"

		err = db.WriteToDB(dbname, data, bucketName)
		assert.NoError(t, err)

		// Verify that the data was written correctly
		db, err := bolt.Open(actualPath, 0o600, nil)
		assert.NoError(t, err)
		defer db.Close()

		err = db.View(func(tx *bolt.Tx) error {
			bucket := tx.Bucket([]byte(bucketName))
			assert.NotNil(t, bucket)

			for k, v := range data {
				value := bucket.Get([]byte(k))
				assert.Equal(t, []byte(v), value)
			}

			return nil
		})
		assert.NoError(t, err)
	})

	t.Run("Test WriteToDB with empty bucketname", func(t *testing.T) {
		actualPath, err := db.GetDBLoc(dbname)
		assert.NoError(t, err)

		t.Cleanup(func() {
			os.Remove(actualPath)
		})

		data := map[string]string{
			"key1": "value1",
			"key2": "value2",
			"key3": "value3",
		}
		bucketName := ""

		err = db.WriteToDB(dbname, data, bucketName)

		require.ErrorIs(t, err, db.ErrCreateBucket)
	})

	t.Run("Test WriteToDB with empty map key", func(t *testing.T) {
		actualPath, err := db.GetDBLoc(dbname)
		assert.NoError(t, err)

		t.Cleanup(func() {
			os.Remove(actualPath)
		})

		data := map[string]string{
			"": "value1",
		}

		err = db.WriteToDB(dbname, data, bucketName)

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

		require.ErrorIs(t, err, db.ErrOpenDB)
	})
}

func Test_DeleteFromDb(t *testing.T) {
	t.Run("Test DeleteFromDb", func(t *testing.T) {

		actualPath, err := db.GetDBLoc(dbname)
		assert.NoError(t, err)

		t.Cleanup(func() {
			os.Remove(actualPath)
		})

		data := map[string]string{
			"key1": "value1",
			"key2": "value2",
			"key3": "value3",
		}
		key := "key1"

		err = db.WriteToDB(dbname, data, bucketName)
		assert.NoError(t, err)

		err = db.DeleteFromDb(dbname, key, bucketName)
		assert.NoError(t, err)

		// Verify that the key was deleted
		db, err := bolt.Open(actualPath, 0o600, nil)
		assert.NoError(t, err)
		defer db.Close()

		err = db.View(func(tx *bolt.Tx) error {
			bucket := tx.Bucket([]byte(bucketName))
			assert.NotNil(t, bucket)

			value := bucket.Get([]byte(key))
			assert.Nil(t, value)

			return nil
		})
		assert.NoError(t, err)
	})

	t.Run("Test DeleteFromDb with empty dbname", func(t *testing.T) {
		dbname := ""
		key := "key1"
		expectedErr := db.ErrDBNameEmpty

		err := db.DeleteFromDb(dbname, key, bucketName)

		require.ErrorIs(t, err, expectedErr)
	})

	t.Run("Test DeleteFromDb with key not found", func(t *testing.T) {

		actualPath, err := db.GetDBLoc(dbname)
		assert.NoError(t, err)

		t.Cleanup(func() {
			os.Remove(actualPath)
		})

		data := map[string]string{
			"key1": "value1",
			"key2": "value2",
			"key3": "value3",
		}
		key := "key4"

		err = db.WriteToDB(dbname, data, bucketName)
		assert.NoError(t, err)

		err = db.DeleteFromDb(dbname, key, bucketName)

		assert.NoError(t, err)
	})

	t.Run("Test DeleteFromDb with bucket not found", func(t *testing.T) {

		key := "key1"
		expectedErr := db.ErrBucketNotFound

		err := db.DeleteFromDb(dbname, key, bucketName)

		require.ErrorIs(t, err, expectedErr)
	})
}

func Test_ListAllRecords(t *testing.T) {
	t.Run("Test ListAllRecords", func(t *testing.T) {

		actualPath, err := db.GetDBLoc(dbname)
		assert.NoError(t, err)

		t.Cleanup(func() {
			os.Remove(actualPath)
		})

		data := map[string]string{
			"key1": "value1",
			"key2": "value2",
			"key3": "value3",
		}

		err = db.WriteToDB(dbname, data, bucketName)
		assert.NoError(t, err)

		records, err := db.GetAllRecords(dbname, bucketName)

		assert.NoError(t, err)
		assert.Equal(t,data, records)
	})

	t.Run("Test ListAllRecords with empty dbname", func(t *testing.T) {
		dbname := ""
		expectedErr := db.ErrDBNameEmpty
		expectedValue := map[string]string{}

		records, err := db.GetAllRecords(dbname, bucketName)

		require.ErrorIs(t, err, expectedErr)
		assert.Equal(t, expectedValue, records)
	})

	t.Run("Test ListAllRecords with bucket not found", func(t *testing.T) {
		expectedErr := db.ErrBucketNotFound

		records, err := db.GetAllRecords(dbname, bucketName)

		require.ErrorIs(t, err, expectedErr)
		assert.Nil(t, records)
	})
}
func Test_UpdateRec(t *testing.T) {
	t.Run("Test UpdateRec", func(t *testing.T) {

		actualPath, err := db.GetDBLoc(dbname)
		assert.NoError(t, err)

		t.Cleanup(func() {
			os.Remove(actualPath)
		})

		data := map[string]string{
			"key1": "value1",
			"key2": "value2",
			"key3": "value3",
		}
		bucketName := "testBucket"
		key := "key1"
		newValue := "updatedValue"

		err = db.WriteToDB(dbname, data, bucketName)
		assert.NoError(t, err)

		err = db.UpdateRec(dbname, key, newValue, bucketName)
		assert.NoError(t, err)

		// Verify that the value was updated
		db, err := bolt.Open(actualPath, 0o600, nil)
		assert.NoError(t, err)
		defer db.Close()

		err = db.View(func(tx *bolt.Tx) error {
			bucket := tx.Bucket([]byte(bucketName))
			assert.NotNil(t, bucket)

			value := bucket.Get([]byte(key))
			assert.Equal(t, []byte(newValue), value)
			return nil
		})
		assert.NoError(t, err)
	})

	t.Run("Test UpdateRec with empty dbname", func(t *testing.T) {
		dbname := ""
		key := "key1"
		newValue := "updatedValue"
		err := db.UpdateRec(dbname, key, newValue, bucketName)

		require.ErrorIs(t, err, db.ErrDBNameEmpty)
	})

	t.Run("Test UpdateRec with key not found", func(t *testing.T) {

		actualPath, err := db.GetDBLoc(dbname)
		assert.NoError(t, err)

		t.Cleanup(func() {
			os.Remove(actualPath)
		})

		data := map[string]string{
			"key1": "value1",
			"key2": "value2",
			"key3": "value3",
		}
		key := "key4"
		newValue := "updatedValue"

		err = db.WriteToDB(dbname, data, bucketName)
		assert.NoError(t, err)

		err = db.UpdateRec(dbname, key, newValue, bucketName)

		require.ErrorIs(t, err, db.ErrProjectNotFound)
	})

	t.Run("Test UpdateRec with bucket not found", func(t *testing.T) {

		key := "key1"
		newValue := "updatedValue"
		expectedErr := db.ErrBucketNotFound

		err := db.UpdateRec(dbname, key, newValue, bucketName)

		require.ErrorIs(t, err, expectedErr)
	})
}

func Test_DeleteDb(t *testing.T) {
	t.Run("Test DeleteDb", func(t *testing.T) {
		actualPath, err := db.GetDBLoc(dbname)
		assert.NoError(t, err)

		t.Cleanup(func() {
			os.Remove(actualPath)
		})

		err = db.DeleteDb(dbname)
		assert.NoError(t, err)

		// Verify that the database file is deleted
		_, err = os.Stat(actualPath)
		assert.True(t, os.IsNotExist(err))
	})

	t.Run("Test DeleteDb with empty dbname", func(t *testing.T) {
		dbname := ""
		expectedErr := db.ErrDBNameEmpty

		err := db.DeleteDb(dbname)

		require.ErrorIs(t, err, expectedErr)
	})
}
