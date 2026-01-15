package pgstorage

import "time"

type bucketNum uint16

type PersonDB struct {
	ID        int64     `db:"id"`
	Name      string    `db:"name"`
	Birthday  time.Time `db:"birthday"`
	Category  string    `db:"category"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type GiftIdeaDB struct {
	ID        int64     `db:"id"`
	PersonID  int64     `db:"person_id"`
	Idea      string    `db:"idea"`
	Notes     string    `db:"notes"`
	CreatedAt time.Time `db:"created_at"`
}

type GiftHistoryDB struct {
	ID        int64     `db:"id"`
	PersonID  int64     `db:"person_id"`
	Gift      string    `db:"gift"`
	Year      int32     `db:"year"`
	CreatedAt time.Time `db:"created_at"`
}

const (
	bucketPrefix     = "bucket_"
	personsTable     = "persons"
	giftIdeasTable   = "gift_ideas"
	giftHistoryTable = "gift_history"
	globalSeqName    = "global_person_id_seq"
)
