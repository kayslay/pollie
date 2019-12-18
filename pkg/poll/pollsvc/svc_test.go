package pollsvc

import (
	"log"
	"pollie/pkg/poll/pollrepo"
	"pollie/pkg/poll/pollrepo/mock"
	"testing"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestPollSvc(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Pollie Service Suite")
}

var _ = Describe("Pkg/Poll/Pollsvc/Svc", func() {

	var (
		pollRepo pollrepo.Repository
		pollSvc  Service
		ctrl     *gomock.Controller
	)

	BeforeEach(func() {
		log.Println("hello")
		ctrl = gomock.NewController(GinkgoT())
	})

	Describe("flow", func() {
		BeforeEach(func() {
			log.Println("Hi")

			pollMock := mock.NewMockRepository(ctrl)
			pollMock.EXPECT().Create(gomock.Any())
		})

		It("holla", func() {
			Expect(nil).To(BeNil())
		})
	})
})
