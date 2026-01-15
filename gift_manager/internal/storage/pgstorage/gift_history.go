package pgstorage

import (
	"context"
	"time"

	"gift_manager/internal/models"

	"github.com/Masterminds/squirrel"
)

func (pg *PGstorage) AddGiftHistory(ctx context.Context, history *models.GiftHistory) error {
	bucket := pg.BucketByPersonID(history.PersonID)

	query := squirrel.Insert(pg.TableWithBucket(bucket, giftHistoryTable)).
		Columns("person_id", "gift", "year", "created_at").
		Values(
			history.PersonID,
			history.Gift,
			history.Year,
			time.Now(),
		).
		PlaceholderFormat(squirrel.Dollar)

	sql, args, err := query.ToSql()
	if err != nil {
		return err
	}

	_, err = pg.db.Exec(ctx, sql, args...)
	return err
}

func (pg *PGstorage) GetGiftHistoryByPersonID(ctx context.Context, personID int64) ([]*models.GiftHistory, error) {
	bucket := pg.BucketByPersonID(personID)

	query := squirrel.Select("id", "person_id", "gift", "year", "created_at").
		From(pg.TableWithBucket(bucket, giftHistoryTable)).
		Where(squirrel.Eq{"person_id": personID}).
		OrderBy("year DESC, created_at DESC").
		PlaceholderFormat(squirrel.Dollar)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := pg.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var history []*models.GiftHistory
	for rows.Next() {
		var item models.GiftHistory
		var createdAt time.Time
		if err := rows.Scan(&item.ID, &item.PersonID, &item.Gift, &item.Year, &createdAt); err != nil {
			return nil, err
		}
		history = append(history, &item)
	}

	return history, nil
}
