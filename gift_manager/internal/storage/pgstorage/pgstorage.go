package pgstorage

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
)

type PGstorage struct {
	db             *pgxpool.Pool
	bucketQuantity uint16
}

func NewPGStorage(connString string, bucketQuantity uint16) (*PGstorage, error) {
	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, errors.Wrap(err, "parse config")
	}

	db, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, errors.Wrap(err, "connect to db")
	}

	if bucketQuantity == 0 {
		return nil, errors.New("bucket quantity must be > 0")
	}

	storage := &PGstorage{
		db:             db,
		bucketQuantity: bucketQuantity,
	}

	if err := storage.initTables(); err != nil {
		return nil, errors.Wrap(err, "init tables")
	}

	return storage, nil
}

func (pg *PGstorage) initTables() error {
	// создаем глобальную sequence
	seqSQL := fmt.Sprintf("CREATE SEQUENCE IF NOT EXISTS %s", globalSeqName)
	if _, err := pg.db.Exec(context.Background(), seqSQL); err != nil {
		return errors.Wrap(err, "create sequence")
	}

	// создаем схемы и таблицы
	for i := 0; i < int(pg.bucketQuantity); i++ {
		schemaName := fmt.Sprintf("%s%d", bucketPrefix, i)

		// создаем схему
		if _, err := pg.db.Exec(context.Background(),
			fmt.Sprintf("CREATE SCHEMA IF NOT EXISTS %s", schemaName)); err != nil {
			return errors.Wrapf(err, "create schema %s", schemaName)
		}

		// таблица persons
		personsSQL := fmt.Sprintf(`
			CREATE TABLE IF NOT EXISTS %s.%s (
				id BIGINT PRIMARY KEY DEFAULT nextval('%s'),
				name VARCHAR(255) NOT NULL,
				birthday DATE NOT NULL,
				category VARCHAR(50),
				created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
				updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
			)`, schemaName, personsTable, globalSeqName)

		// таблица gift_ideas
		giftIdeasSQL := fmt.Sprintf(`
			CREATE TABLE IF NOT EXISTS %s.%s (
				id BIGSERIAL PRIMARY KEY,
				person_id BIGINT NOT NULL,
				idea TEXT NOT NULL,
				notes TEXT,
				created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
				FOREIGN KEY (person_id) REFERENCES %s.%s(id) ON DELETE CASCADE
			)`, schemaName, giftIdeasTable, schemaName, personsTable)

		// таблица gift_history
		giftHistorySQL := fmt.Sprintf(`
			CREATE TABLE IF NOT EXISTS %s.%s (
				id BIGSERIAL PRIMARY KEY,
				person_id BIGINT NOT NULL,
				gift TEXT NOT NULL,
				year INTEGER NOT NULL,
				created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
				FOREIGN KEY (person_id) REFERENCES %s.%s(id) ON DELETE CASCADE
			)`, schemaName, giftHistoryTable, schemaName, personsTable)

		// выполняем sql
		queries := []string{personsSQL, giftIdeasSQL, giftHistorySQL}
		for _, query := range queries {
			if _, err := pg.db.Exec(context.Background(), query); err != nil {
				return errors.Wrapf(err, "create table in schema %s", schemaName)
			}
		}

		// создаем индексы
		indexes := []string{
			fmt.Sprintf("CREATE INDEX IF NOT EXISTS idx_%s_person ON %s.%s(person_id)",
				giftIdeasTable, schemaName, giftIdeasTable),
			fmt.Sprintf("CREATE INDEX IF NOT EXISTS idx_%s_person ON %s.%s(person_id)",
				giftHistoryTable, schemaName, giftHistoryTable),
		}

		for _, index := range indexes {
			pg.db.Exec(context.Background(), index)
		}
	}
	return nil
}

func (pg *PGstorage) Close() {
	if pg.db != nil {
		pg.db.Close()
	}
}
