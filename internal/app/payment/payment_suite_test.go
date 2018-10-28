package payment_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestPayment(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Payment Suite")
}
