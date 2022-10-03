package db

import (
	"github.com/boltdb/bolt"
	"github.com/rockjoon/nomadcoin/utils"
)

const (
	dbName     = "blockchain.db"
	checkpoint = "checkpoint"
)

var db *bolt.DB

func DB() *bolt.DB {
	if db == nil {
		dbPointer, err := bolt.Open(dbName, 0600, nil)
		utils.HandleError(err)
		db = dbPointer
		err = db.Update(func(tx *bolt.Tx) error {
			_, bucketErr := tx.CreateBucketIfNotExists([]byte(DataBucket))
			utils.HandleError(bucketErr)
			_, bucketErr = tx.CreateBucketIfNotExists([]byte(BlocksBucket))
			utils.HandleError(bucketErr)
			return nil
		})
		utils.HandleError(err)
	}
	return db
}

func Close() {
	DB().Close()
}

func Save(bucket Bucket, key string, value interface{}) {
	err := DB().Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucket))
		return bucket.Put([]byte(key), utils.ToBytes(value))
	})
	utils.HandleError(err)

}

func SaveBlock(key string, value interface{}) {
	Save(BlocksBucket, key, value)
}

func SaveBlockChain(value interface{}) {
	Save(DataBucket, checkpoint, value)
}

func GetCheckPoint() []byte {
	var data []byte
	DB().View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(DataBucket))
		data = bucket.Get([]byte(checkpoint))
		return nil
	})
	return data
}

func Block(hash string) []byte {
	var data []byte
	DB().View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BlocksBucket))
		data = bucket.Get([]byte(hash))
		return nil
	})
	return data
}
