package exercise2

import (
	"fmt"

	"github.com/boltdb/bolt"
)

//InitBoltDB will create the needed data for our bolt db
func InitBoltDB() error {
	// Open the my.db data file in your current directory.
	// It will be created if it doesn't exist.
	db, err := bolt.Open("bolturls.db", 0600, nil)
	if err != nil {
		return fmt.Errorf("Error opening bolt db: %v", err)
	}
	defer db.Close()

	err = db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("UrlBucket"))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		err = b.Put([]byte("/bolt-github"), []byte("https://github.com/boltdb/bolt"))
		if err != nil {
			return fmt.Errorf("put bucket: %s", err)
		}
		err = b.Put([]byte("/bolt-godoc"), []byte("https://godoc.org/github.com/boltdb/bolt"))
		if err != nil {
			return fmt.Errorf("put bucket: %s", err)
		}
		return nil
	})
	return err
}

func getBoltRecords() (map[string]string, error) {
	db, err := bolt.Open("bolturls.db", 0600, nil)
	if err != nil {
		return nil, fmt.Errorf("Error opening bolt db: %v", err)
	}
	defer db.Close()
	urls := make(map[string]string)
	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("UrlBucket"))
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			urls[string(k)] = string(v)
		}
		return nil
	})
	return urls, err
}
