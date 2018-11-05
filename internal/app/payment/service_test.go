package payment_test

import (
	"context"
	"database/sql"
	"encoding/json"
	"github.com/carlosroman/payments-api/internal/app/payment"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

var _ = Describe("Service", func() {

	var (
		s      payment.Service
		db     *sql.DB
		dbMock sqlmock.Sqlmock
		ctx    context.Context
	)

	BeforeEach(func() {
		d, mock, err := sqlmock.New()
		Expect(err).ShouldNot(HaveOccurred())
		db = d
		s = payment.NewService(db)
		dbMock = mock
		ctx = context.Background()
	})

	Describe("Saving a new payment", func() {
		Context("when successful", func() {
			It("should return id from DB", func() {
				expectedId := "expectedId"
				rows := sqlmock.NewRows([]string{"ID"}).AddRow(expectedId)
				dbMock.ExpectQuery("INSERT INTO payments\\(ID, info\\)").
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnRows(rows)

				p := payment.Payment{Attributes: payment.Attributes{Reference: "some ref"}}
				id, err := s.Save(ctx, p)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(id).Should(Equal(expectedId))
			})

			It("should save correct object", func() {
				expectedId := "Some Id"
				s = payment.NewServiceWithUuidGen(db, func() string {
					return expectedId
				})
				p := payment.Payment{Attributes: payment.Attributes{Reference: "some ref"}, Id: expectedId}
				bs, err := json.Marshal(p)
				Expect(err).ShouldNot(HaveOccurred())
				rows := sqlmock.NewRows([]string{"ID"}).AddRow("some id")
				dbMock.ExpectQuery("INSERT INTO payments\\(ID, info\\)").
					WithArgs(sqlmock.AnyArg(), string(bs)).
					WillReturnRows(rows)
				_, err = s.Save(ctx, payment.Payment{Attributes: payment.Attributes{Reference: "some ref"}})
				Expect(err).ShouldNot(HaveOccurred())
				Expect(dbMock.ExpectationsWereMet()).ShouldNot(HaveOccurred())
			})
		})
	})

	Describe("Getting a payment", func() {
		Context("when successful", func() {
			It("should return payment from DB", func() {
				p := payment.Payment{Attributes: payment.Attributes{Reference: "some ref"}, Id: "awesome id"}
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

			It("should return empty slice", func() {
				rows := sqlmock.NewRows([]string{"info"})
				dbMock.ExpectQuery("SELECT info FROM payments WHERE info ->> 'organisation_id' = ?").
					WithArgs("OrgId").
					WillReturnRows(rows)
				actual, err := s.SearchByOrganisationId(ctx, "OrgId")
				Expect(err).ShouldNot(HaveOccurred())
				Expect(actual).ShouldNot(BeNil())
			})
		})

		Context("when not successful", func() {
			It("should return db error back", func() {
				dbMock.ExpectQuery("SELECT info FROM payments WHERE info ->> 'organisation_id' = ?").
					WithArgs("OrgId").
					WillReturnError(sql.ErrConnDone)
				_, err := s.SearchByOrganisationId(ctx, "OrgId")
				Expect(err).Should(HaveOccurred())
				Expect(err).To(Equal(sql.ErrConnDone))
			})
		})
	})

	Describe("when HealthCheck called", func() {
		var mockDb mockDatabase
		BeforeEach(func() {
			mockDb = mockDatabase{}
			s = payment.NewService(&mockDb)
		})

		Context("when healthy", func() {
			It("should return healthy status", func() {
				mockDb.On("PingContext", mock.Anything).
					Return(nil)

				actual := s.HealthCheck(ctx)
				expected := payment.HealthCheckStatus{
					Healthy: true,
					Message: "okay",
				}
				Expect(actual).Should(Equal(expected))
				mockDb.AssertCalled(GinkgoT(), "PingContext", ctx)
			})
		})

		Context("when unhealthy", func() {
			It("should return unhealthy status", func() {
				mockDb.On("PingContext", mock.Anything).
					Return(sql.ErrConnDone)

				actual := s.HealthCheck(ctx)
				expected := payment.HealthCheckStatus{
					Healthy: false,
					Message: "sql: connection is already closed",
				}
				Expect(actual).Should(Equal(expected))
				mockDb.AssertCalled(GinkgoT(), "PingContext", ctx)
			})
		})
	})
})

type mockDatabase struct {
	mock.Mock
}

func (m *mockDatabase) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	arg := m.Called(ctx, query, args)
	return arg.Get(0).(*sql.Row)
}

func (m *mockDatabase) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	arg := m.Called(ctx, query, args)
	return arg.Get(0).(*sql.Rows), arg.Error(1)

}
func (m *mockDatabase) PingContext(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}
