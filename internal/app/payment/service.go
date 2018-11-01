package payment

import (
	"context"
	"database/sql"
	"encoding/json"
	log "github.com/sirupsen/logrus"
)

type Service interface {
	Save(ctx context.Context, payment Payment) (id string, err error)
}

type Database interface {
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
}

func NewService(db Database) Service {
	return &service{
		db: db,
	}
}

type service struct {
	db Database
}

func (s *service) Save(ctx context.Context, payment Payment) (id string, err error) {
	bs, err := json.Marshal(payment)
	if err != nil {
		return id, err
	}

	err = s.db.QueryRowContext(ctx,
		"INSERT INTO payments(info) VALUES(?) returning ID",
		string(bs)).Scan(&id)
	if err != nil {
		return id, err
	}

	log.Infof("Inserted payment, id is '%s'", id)
	return id, err
}
