package pollie_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestPollie(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Pollie Suite")
}
