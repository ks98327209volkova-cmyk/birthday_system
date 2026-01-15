package pgstorage

import (
	"context"
	"time"

	"gift_manager/internal/models"

	"github.com/Masterminds/squirrel"
)

func (pg *PGstorage) AddGiftIdea(ctx context.Context, idea *models.GiftIdea) error {
	bucket := pg.BucketByPersonID(idea.PersonID)

	query := squirrel.Insert(pg.TableWithBucket(bucket, giftIdeasTable)).
		Columns("person_id", "idea", "notes", "created_at").
		Values(
			idea.PersonID,
			idea.Idea,
			idea.Notes,
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

func (pg *PGstorage) GetGiftIdeasByPersonID(ctx context.Context, personID int64) ([]*models.GiftIdea, error) {
	bucket := pg.BucketByPersonID(personID)

	query := squirrel.Select("id", "person_id", "idea", "notes", "created_at").
		From(pg.TableWithBucket(bucket, giftIdeasTable)).
		Where(squirrel.Eq{"person_id": personID}).
		OrderBy("created_at DESC").
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

	var ideas []*models.GiftIdea
	for rows.Next() {
		var idea models.GiftIdea
		var createdAt time.Time
		if err := rows.Scan(&idea.ID, &idea.PersonID, &idea.Idea, &idea.Notes, &createdAt); err != nil {
			return nil, err
		}
		ideas = append(ideas, &idea)
	}

	return ideas, nil
}

func (pg *PGstorage) DeleteGiftIdea(ctx context.Context, personID, ideaID int64) error {
	bucket := pg.BucketByPersonID(personID)

	query := squirrel.Delete(pg.TableWithBucket(bucket, giftIdeasTable)).
		Where(squirrel.Eq{"id": ideaID}).
		PlaceholderFormat(squirrel.Dollar)

	sql, args, err := query.ToSql()
	if err != nil {
		return err
	}

	_, err = pg.db.Exec(ctx, sql, args...)
	return err
}
