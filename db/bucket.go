package db

type Bucket string

const (
	DataBucket   = Bucket("data")
	BlocksBucket = Bucket("blocks")
)
