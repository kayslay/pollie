package voterepo

import (
	"context"
	"fmt"
	"pollie/config"
	"pollie/models"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestVote(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Pollie Suite")
}

var (
	ipExists   = "127.0.0.2"
	ipNoExists = "127.0.0.1"
)

var _ = Describe("Pkg/Vote/Voterepo/Repo", func() {

	viper.AddConfigPath("./../../../")
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	viper.Set("MGO_DB", "pollie_tests")

	var (
		repo   Repository
		err    error
		mgoFn  config.MgoFn
		id     primitive.ObjectID = primitive.NewObjectID()
		pollID primitive.ObjectID = primitive.NewObjectID()
	)

	BeforeEach(func() {
		mgoFn, err = config.NewMgoDB()
		if err != nil {
			Fail("could not get db")
		}

		setUpTests(mgoFn, id, pollID)
		repo = NewRepository(mgoFn)
	})

	AfterEach(func() {
		mgoFn(models.VoteCollection).Drop(context.TODO())
		mgoFn(models.PollCollection).Drop(context.TODO())
		mgoFn("ips").Drop(context.TODO())
	})

	Describe("Check if IPExist", func() {

		It("expect ip to exists", func() {
			b, err := repo.IPExists(pollID.Hex(), ipExists)
			Expect(err).To(BeNil())
			Expect(b).To(Equal(true))
		})

		It("expect ip not to exists", func() {
			b, err := repo.IPExists(pollID.Hex(), ipNoExists)
			Expect(err).To(BeNil())
			Expect(b).To(Equal(false))
		})
	})

	Describe("Vote for a poll", func() {

		It("option passed should pass", func() {
			id, err := repo.Vote(
				models.Vote{
					ID:     primitive.NewObjectID(),
					PollID: pollID,
					Option: []int{5},
					Meta: models.Meta{
						IP: ipExists,
					},
				},
			)
			fmt.Println(id)
			Expect(err).To(BeNil())

		})
	})

})

func setUpTests(mgoFn config.MgoFn, id, pollID primitive.ObjectID) {
	n, err := mgoFn(models.PollCollection).InsertMany(context.TODO(),
		[]interface{}{
			models.Poll{
				ID:          pollID,
				Type:        models.PTypeSingle,
				Description: "test",
				ShortCode:   "Nhyeia",
				Option:      []string{"a", "b"},
				Summary: models.PollSummary{
					OptionCount: make([]int64, 2),
				},
			},
		},
	)

	Expect(err).To(BeNil())
	Expect(len(n.InsertedIDs)).To(Equal(1))
	// set up vote
	n, err = mgoFn(models.VoteCollection).InsertMany(context.TODO(),
		[]interface{}{
			models.Vote{
				ID:     id,
				PollID: pollID,
				Option: []int{},
				Meta: models.Meta{
					IP: ipExists,
				},
			},
		},
	)
	Expect(err).To(BeNil())
	Expect(len(n.InsertedIDs)).To(Equal(1))

}
