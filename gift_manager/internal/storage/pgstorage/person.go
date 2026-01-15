package pgstorage

import (
	"context"
	"fmt"
	"time"

	"gift_manager/internal/models"

	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
)

func (pg *PGstorage) AddPerson(ctx context.Context, person *models.Person) (int64, error) {
	var id int64
	err := pg.db.QueryRow(ctx, fmt.Sprintf("SELECT nextval('%s')", globalSeqName)).Scan(&id)
	if err != nil {
		return 0, errors.Wrap(err, "get next id")
	}

	bucket := pg.BucketByPersonID(id)

	query := squirrel.Insert(pg.TableWithBucket(bucket, personsTable)).
		Columns("id", "name", "birthday", "category", "created_at", "updated_at").
		Values(
			id,
			person.Name,
			person.Birthday,
			string(person.Category),
			time.Now(),
			time.Now(),
		).
		PlaceholderFormat(squirrel.Dollar)

	sql, args, err := query.ToSql()
	if err != nil {
		return 0, err
	}

	_, err = pg.db.Exec(ctx, sql, args...)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (pg *PGstorage) GetPersonByID(ctx context.Context, id int64) (*models.Person, error) {
	bucket := pg.BucketByPersonID(id)

	query := squirrel.Select("id", "name", "birthday", "category").
		From(pg.TableWithBucket(bucket, personsTable)).
		Where(squirrel.Eq{"id": id}).
		PlaceholderFormat(squirrel.Dollar)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	var person models.Person
	var birthday time.Time
	err = pg.db.QueryRow(ctx, sql, args...).Scan(
		&person.ID, &person.Name, &birthday, &person.Category,
	)

	if err != nil {
		return nil, errors.Wrap(err, "person not found")
	}

	person.Birthday = birthday.Format("2006-01-02")
	return &person, nil
}

func (pg *PGstorage) GetAllPersons(ctx context.Context) ([]*models.Person, error) {
	var all []*models.Person

	for i := 0; i < int(pg.bucketQuantity); i++ {
		bucket := bucketNum(i)

		query := squirrel.Select("id", "name", "birthday", "category").
			From(pg.TableWithBucket(bucket, personsTable)).
			PlaceholderFormat(squirrel.Dollar)

		sql, args, err := query.ToSql()
		if err != nil {
			return nil, err
		}

		rows, err := pg.db.Query(ctx, sql, args...)
		if err != nil {
			return nil, err
		}

		for rows.Next() {
			var p models.Person
			var bd time.Time
			if err := rows.Scan(&p.ID, &p.Name, &bd, &p.Category); err != nil {
				rows.Close()
				return nil, err
			}
			p.Birthday = bd.Format("2006-01-02")
			all = append(all, &p)
		}
		rows.Close()
	}

	return all, nil
}

func (pg *PGstorage) UpdatePerson(ctx context.Context, person *models.Person) error {
	bucket := pg.BucketByPersonID(person.ID)

	query := squirrel.Update(pg.TableWithBucket(bucket, personsTable)).
		Set("name", person.Name).
		Set("birthday", person.Birthday).
		Set("category", string(person.Category)).
		Set("updated_at", time.Now()).
		Where(squirrel.Eq{"id": person.ID}).
		PlaceholderFormat(squirrel.Dollar)

	sql, args, err := query.ToSql()
	if err != nil {
		return err
	}

	_, err = pg.db.Exec(ctx, sql, args...)
	return err
}

func (pg *PGstorage) DeletePerson(ctx context.Context, id int64) error {
	bucket := pg.BucketByPersonID(id)

	query := squirrel.Delete(pg.TableWithBucket(bucket, personsTable)).
		Where(squirrel.Eq{"id": id}).
		PlaceholderFormat(squirrel.Dollar)

	sql, args, err := query.ToSql()
	if err != nil {
		return err
	}

	_, err = pg.db.Exec(ctx, sql, args...)
	return err
}
