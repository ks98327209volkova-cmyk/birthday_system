package pgstorage

import "fmt"

func (pg *PGstorage) BucketByPersonID(personID int64) bucketNum {
	if pg.bucketQuantity == 0 {
		return 0
	}
	return bucketNum(personID % int64(pg.bucketQuantity))
}

func (pg *PGstorage) TableWithBucket(bucket bucketNum, tableName string) string {
	return fmt.Sprintf("%s%d.%s", bucketPrefix, bucket, tableName)
}
