package pollrepo

import (
	"context"
	"pollie/config"
	"pollie/models"
	"testing"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestPollRepo(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Pollie Repository Suite")
}

var _ = Describe("Pkg/Poll/Pollrepo/Repo", func() {
	viper.AddConfigPath("./../../../")
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	viper.Set("MGO_DB", "pollie_tests")

	var (
		repo      Repository
		err       error
		mgoFn     config.MgoFn
		passID    primitive.ObjectID = primitive.NewObjectID()
		deletedID primitive.ObjectID = primitive.NewObjectID()
	)

	BeforeEach(func() {
		mgoFn, err = config.NewMgoDB()
		if err != nil {
			Fail("could not get db")
		}
		repo = NewRepository(mgoFn)
	})

	AfterEach(func() {
		mgoFn(models.VoteCollection).Drop(context.TODO())
		mgoFn(models.PollCollection).Drop(context.TODO())
		mgoFn("ips").Drop(context.TODO())
	})

	Describe("Get poll", func() {
		BeforeEach(func() {
			expireAt := time.Now().Add(time.Minute)
			n, err := mgoFn(models.PollCollection).InsertMany(context.TODO(),
				[]interface{}{
					models.Poll{
						ID:          passID,
						Type:        models.PTypeSingle,
						Description: "test",
						ShortCode:   "PassSht",
						Option:      []string{"a", "b"},
						Summary: models.PollSummary{
							OptionCount: make([]int64, 2),
						},
						ExpiresAt: &expireAt,
					},
				},
			)

			Expect(err).To(BeNil())
			Expect(len(n.InsertedIDs)).To(Equal(1))

			// expired poll
			deletedAt := time.Now().Add(-1 * time.Minute)
			n, err = mgoFn(models.PollCollection).InsertMany(context.TODO(),
				[]interface{}{
					models.Poll{
						ID:          deletedID,
						Type:        models.PTypeSingle,
						Description: "test",
						ShortCode:   "Deletd",
						Option:      []string{"a", "b"},
						Summary: models.PollSummary{
							OptionCount: make([]int64, 2),
						},
						DeletedAt: &deletedAt,
					},
				},
			)

			Expect(err).To(BeNil())
			Expect(len(n.InsertedIDs)).To(Equal(1))
		})

		Context("Get Valid Poll", func() {
			It("Get Poll using ID", func() {
				_, err := repo.Get(context.TODO(), passID.Hex())
				Expect(err).To(BeNil())
			})

			It("Get Poll using Shortcode", func() {
				_, err := repo.Get(context.TODO(), "PassSht")
				Expect(err).To(BeNil())
			})
		})

		Context("Should return Error when trying to get deleted poll", func() {
			It("Get with ID", func() {
				_, err := repo.Get(context.TODO(), deletedID.Hex())
				Expect(err).NotTo(BeNil())
			})

			It("Get with shortCode", func() {
				_, err := repo.Get(context.TODO(), "Deletd")
				Expect(err).NotTo(BeNil())
			})
		})

	})
})
