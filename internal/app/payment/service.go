package payment

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	log "github.com/sirupsen/logrus"
)

type Service interface {
	Save(ctx context.Context, payment Payment) (id string, err error)
	Get(ctx context.Context, paymentId string) (payment Payment, err error)
}

var ErrNotFound = errors.New("payment: not found")

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
		"INSERT INTO payments(info) VALUES($1) returning ID;",
		string(bs)).Scan(&id)
	if err != nil {
		return id, err
	}

	log.Infof("Inserted payment, id is '%s'", id)
	return id, err
}

func (s *service) Get(ctx context.Context, paymentId string) (payment Payment, err error) {
	var info string
	err = s.db.QueryRowContext(ctx,
		"SELECT info FROM payments WHERE ID = $1;",
		paymentId).Scan(&info)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return payment, ErrNotFound
		default:
			return payment, err
		}
	}
	err = json.Unmarshal([]byte(info), &payment)
	return payment, err
}
