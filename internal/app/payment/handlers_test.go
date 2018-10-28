package payment_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/carlosroman/payments-api/internal/app/payment"
	"github.com/gorilla/mux"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/mock"
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

				ms.On("Save", mock.AnythingOfType("payment.Payment")).Return("new-payment-id", nil)

				resp, err := http.DefaultClient.Do(req)
				Expect(err).ShouldNot(HaveOccurred())
				defer resp.Body.Close()
				Expect(resp.StatusCode).To(Equal(http.StatusCreated))
				ms.AssertExpectations(GinkgoT())
			})

			It("should return location", func() {
				req := givenValidPaymentRequest(ts.URL)

				ms.On("Save", mock.AnythingOfType("payment.Payment")).Return("new-payment-id", nil)

				resp, err := http.DefaultClient.Do(req)
				Expect(err).ShouldNot(HaveOccurred())
				defer resp.Body.Close()
				Expect(resp.Header.Get("Location")).To(Equal("/new-payment-id"))
				ms.AssertExpectations(GinkgoT())
			})

			It("should call save correctly", func() {
				req := givenValidPaymentRequest(ts.URL)

				ms.On("Save", mock.AnythingOfType("payment.Payment")).Return("new-payment-id", nil)

				resp, err := http.DefaultClient.Do(req)
				Expect(err).ShouldNot(HaveOccurred())
				defer resp.Body.Close()
				ms.AssertCalled(GinkgoT(), "Save", payment.Payment{})
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

				ms.On("Save", mock.AnythingOfType("payment.Payment")).Return("", errors.New("something went wrong"))
				resp, err := http.DefaultClient.Do(req)
				Expect(err).ShouldNot(HaveOccurred())
				defer resp.Body.Close()
				Expect(resp.StatusCode).To(Equal(http.StatusInternalServerError))
			})
		})
	})
})

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

func (s *mockService) Save(payment payment.Payment) (id string, err error) {
	args := s.Called(payment)
	return args.String(0), args.Error(1)
}
