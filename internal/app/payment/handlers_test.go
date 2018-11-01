package payment_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/carlosroman/payments-api/internal/app/payment"
	"github.com/gorilla/mux"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/mock"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
)

func init() {
	log.SetLevel(log.DebugLevel)
}

var _ = Describe("Handlers", func() {

	var (
		ms mockService
		r  *mux.Router
		ts *httptest.Server
	)

	BeforeEach(func() {
		ms = mockService{}
		r = payment.GetHandlers(&ms)
		ts = httptest.NewServer(r)
	})

	AfterEach(func() {
		ts.Close()
	})

	Describe("Saving a new payment", func() {
		Context("that is valid", func() {
			It("should return created", func() {
				req := givenValidPaymentRequest(ts.URL)

				ms.On("Save", mock.AnythingOfType("*context.valueCtx"), mock.AnythingOfType("payment.Payment")).Return("new-payment-id", nil)

				resp, err := http.DefaultClient.Do(req)
				Expect(err).ShouldNot(HaveOccurred())
				defer resp.Body.Close()
				Expect(resp.StatusCode).To(Equal(http.StatusCreated))
				ms.AssertExpectations(GinkgoT())
			})

			It("should return location", func() {
				req := givenValidPaymentRequest(ts.URL)

				ms.On("Save", mock.AnythingOfType("*context.valueCtx"), mock.AnythingOfType("payment.Payment")).Return("new-payment-id", nil)

				resp, err := http.DefaultClient.Do(req)
				Expect(err).ShouldNot(HaveOccurred())
				defer resp.Body.Close()
				Expect(resp.Header.Get("Location")).To(Equal("/payment/new-payment-id"))
				ms.AssertExpectations(GinkgoT())
			})

			It("should call save correctly", func() {
				req := givenValidPaymentRequest(ts.URL)

				ms.On("Save", mock.AnythingOfType("*context.valueCtx"), mock.AnythingOfType("payment.Payment")).Return("new-payment-id", nil)

				resp, err := http.DefaultClient.Do(req)
				Expect(err).ShouldNot(HaveOccurred())
				defer resp.Body.Close()
				ms.AssertCalled(GinkgoT(), "Save", mock.AnythingOfType("*context.valueCtx"), payment.Payment{})
			})
		})

		Context("that is invalid", func() {
			It("should return bad request", func() {
				body := strings.NewReader("this is not valid json \n{}\n")
				req, err := http.NewRequest("POST", fmt.Sprintf("%s/payment", ts.URL), body)
				Expect(err).ShouldNot(HaveOccurred())
				resp, err := http.DefaultClient.Do(req)
				Expect(err).ShouldNot(HaveOccurred())
				defer resp.Body.Close()
				Expect(resp.StatusCode).To(Equal(http.StatusBadRequest))
			})
		})

		Context("when the back end fails", func() {
			It("should return internal server error", func() {
				req := givenValidPaymentRequest(ts.URL)

				ms.On("Save", mock.AnythingOfType("*context.valueCtx"), mock.AnythingOfType("payment.Payment")).Return("", errors.New("something went wrong"))
				resp, err := http.DefaultClient.Do(req)
				Expect(err).ShouldNot(HaveOccurred())
				defer resp.Body.Close()
				Expect(resp.StatusCode).To(Equal(http.StatusInternalServerError))
			})
		})
	})

	Describe("Getting a payment", func() {

		Context("that exists in the db", func() {
			It("should return the payment", func() {

				id, req := givenPaymentRequest(ts.URL)

				expected := payment.Payment{Id: id, Reference: "some ref"}
				ms.On("Get", mock.AnythingOfType("*context.valueCtx"), mock.AnythingOfType("string")).
					Return(expected, nil)

				resp, err := http.DefaultClient.Do(req)
				Expect(err).ShouldNot(HaveOccurred())
				defer resp.Body.Close()

				body, err := ioutil.ReadAll(resp.Body)
				Expect(err).ShouldNot(HaveOccurred())
				var actual payment.Payment
				err = json.Unmarshal(body, &actual)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(actual).Should(Equal(expected))
				ms.AssertCalled(GinkgoT(), "Get", mock.AnythingOfType("*context.valueCtx"), id)
			})

			It("should return 200 status", func() {

				id, req := givenPaymentRequest(ts.URL)

				expected := payment.Payment{Id: id, Reference: "some ref"}
				ms.On("Get", mock.AnythingOfType("*context.valueCtx"), mock.AnythingOfType("string")).
					Return(expected, nil)

				resp, err := http.DefaultClient.Do(req)
				Expect(err).ShouldNot(HaveOccurred())
				defer resp.Body.Close()
				Expect(resp.StatusCode).Should(Equal(200))
			})
		})

		Context("that does not exists in the db", func() {
			It("should return not found", func() {
				_, req := givenPaymentRequest(ts.URL)
				ms.On("Get", mock.AnythingOfType("*context.valueCtx"), mock.AnythingOfType("string")).
					Return(payment.Payment{}, payment.ErrNotFound)
				resp, err := http.DefaultClient.Do(req)
				Expect(err).ShouldNot(HaveOccurred())
				defer resp.Body.Close()
				Expect(resp.StatusCode).To(Equal(http.StatusNotFound))
			})
		})

		Context("when the database errors", func() {
			It("should return internal server error", func() {
				_, req := givenPaymentRequest(ts.URL)
				ms.On("Get", mock.AnythingOfType("*context.valueCtx"), mock.AnythingOfType("string")).
					Return(payment.Payment{}, errors.New("some DB error"))
				resp, err := http.DefaultClient.Do(req)
				Expect(err).ShouldNot(HaveOccurred())
				defer resp.Body.Close()
				Expect(resp.StatusCode).To(Equal(http.StatusInternalServerError))
			})
		})
	})
})

func givenPaymentRequest(url string) (id string, req *http.Request) {
	id = uuid.NewV4().String()
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/payment/%s", url, id), nil)
	Expect(err).ShouldNot(HaveOccurred())
	return id, req
}

func givenValidPaymentRequest(url string) *http.Request {
	p := payment.Payment{}
	bs, err := json.Marshal(p)
	Expect(err).ShouldNot(HaveOccurred())
	body := bytes.NewReader(bs)
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/payment", url), body)
	Expect(err).ShouldNot(HaveOccurred())
	req.Header.Set("Content-Type", "application/json")
	return req
}

type mockService struct {
	mock.Mock
}

func (s *mockService) Save(ctx context.Context, payment payment.Payment) (id string, err error) {
	args := s.Called(ctx, payment)
	return args.String(0), args.Error(1)
}

func (s *mockService) Get(ctx context.Context, id string) (p payment.Payment, err error) {
	args := s.Called(ctx, id)
	return args.Get(0).(payment.Payment), args.Error(1)
}
