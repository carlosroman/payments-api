package payment_test

import (
	"context"
	"encoding/json"
	"github.com/carlosroman/payments-api/internal/app/payment"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

var _ = Describe("Service", func() {

	var (
		s      payment.Service
		dbMock sqlmock.Sqlmock
		ctx    context.Context
	)

	BeforeEach(func() {
		db, mock, err := sqlmock.New()
		Expect(err).ShouldNot(HaveOccurred())
		s = payment.NewService(db)
		dbMock = mock
		ctx = context.Background()
	})

	Describe("Saving a new payment", func() {
		Context("when successful", func() {
			It("should return id from DB", func() {

				expectedId := "expectedId"
				rows := sqlmock.NewRows([]string{"ID"}).AddRow(expectedId)
				dbMock.ExpectQuery("INSERT INTO payments\\(info\\)").
					WithArgs(sqlmock.AnyArg()).
					WillReturnRows(rows)

				p := payment.Payment{Reference: "some ref"}
				id, err := s.Save(ctx, p)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(id).Should(Equal(expectedId))
			})

			It("should save correct object", func() {
				p := payment.Payment{Reference: "some ref"}
				bs, err := json.Marshal(p)
				Expect(err).ShouldNot(HaveOccurred())
				rows := sqlmock.NewRows([]string{"ID"}).AddRow("some id")
				dbMock.ExpectQuery("INSERT INTO payments\\(info\\)").
					WithArgs(string(bs)).
					WillReturnRows(rows)
				_, err = s.Save(ctx, p)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(dbMock.ExpectationsWereMet()).ShouldNot(HaveOccurred())
			})
		})
	})
})
