package payment_test

import (
	"context"
	"database/sql"
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

	Describe("Getting a payment", func() {
		Context("when successful", func() {
			It("should return payment from DB", func() {
				p := payment.Payment{Reference: "some ref", Id: "awesome id"}
				bs, err := json.Marshal(p)
				Expect(err).ShouldNot(HaveOccurred())

				rows := sqlmock.NewRows([]string{"info"}).AddRow(string(bs))
				dbMock.ExpectQuery("SELECT info FROM payments WHERE ID = ?").
					WithArgs(p.Id).
					WillReturnRows(rows)

				actual, err := s.Get(ctx, p.Id)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(actual).To(Equal(p))
			})
		})
		Context("when not successful", func() {
			It("should return not found if no record", func() {

				dbMock.ExpectQuery("SELECT info FROM payments WHERE ID = ?").
					WithArgs("some id").
					WillReturnError(sql.ErrNoRows)
				_, err := s.Get(ctx, "some id")
				Expect(err).Should(HaveOccurred())
				Expect(err).To(Equal(payment.ErrNotFound))
			})
			It("should return all other errors", func() {

				dbMock.ExpectQuery("SELECT info FROM payments WHERE ID = ?").
					WithArgs("some id").
					WillReturnError(sql.ErrConnDone)
				_, err := s.Get(ctx, "some id")
				Expect(err).Should(HaveOccurred())
				Expect(err).To(Equal(sql.ErrConnDone))
			})
		})
	})

	Describe("Search By Organisation Id", func() {
		Context("when successful", func() {
			It("should return all payments for given Organisation Id", func() {
				ps := []payment.Payment{
					{Id: "A", OrganisationId: "OrgId"},
					{Id: "B", OrganisationId: "OrgId"},
					{Id: "C", OrganisationId: "OrgId"},
				}

				rows := sqlmock.NewRows([]string{"info"})

				for _, p := range ps {
					bs, err := json.Marshal(p)
					Expect(err).ShouldNot(HaveOccurred())
					rows.AddRow(string(bs))
				}

				dbMock.ExpectQuery("SELECT info FROM payments WHERE info ->> 'organisation_id' = ?").
					WithArgs("OrgId").
					WillReturnRows(rows)

				actual, err := s.SearchByOrganisationId(ctx, "OrgId")
				Expect(err).ShouldNot(HaveOccurred())
				Expect(actual).To(Equal(ps))
			})
		})
	})
})
