package payment

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	log "github.com/sirupsen/logrus"
)

type HealthCheckStatus struct {
	Message string `json:"message"`
	Healthy bool   `json:"healthy"`
}

type Service interface {
	Save(ctx context.Context, payment Payment) (id string, err error)
	Get(ctx context.Context, paymentId string) (payment Payment, err error)
	SearchByOrganisationId(ctx context.Context, organisationId string) (payments []Payment, err error)
	HealthCheck(ctx context.Context) HealthCheckStatus
}

var ErrNotFound = errors.New("payment: not found")

type Database interface {
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	PingContext(ctx context.Context) error
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

func (s *service) SearchByOrganisationId(ctx context.Context, organisationId string) (payments []Payment, err error) {

	rows, err := s.db.QueryContext(ctx,
		"SELECT info FROM payments WHERE info ->> 'organisation_id' = $1;",
		organisationId)

	if err != nil {
		return payments, err
	}

	payments = make([]Payment, 0)
	for rows.Next() {
		var payment Payment
		var info string
		if err = rows.Scan(&info); err != nil {
			return payments, err
		}
		if err = json.Unmarshal([]byte(info), &payment); err != nil {
			return payments, err
		}
		payments = append(payments, payment)
	}

	return payments, err
}

func (s *service) HealthCheck(ctx context.Context) HealthCheckStatus {
	if err := s.db.PingContext(ctx); err != nil {
		return HealthCheckStatus{
			Message: err.Error(),
			Healthy: false,
		}
	}

	return HealthCheckStatus{
		Message: "okay",
		Healthy: true,
	}
}
