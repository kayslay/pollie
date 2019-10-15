package pollsvc

import (
	"context"
	"pollie/models"
	"pollie/pkg/poll/pollrepo"
	"time"

	str "github.com/kayslay/random-str"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Service defines the behavior of an poll service
type Service interface {
	Create(ctx context.Context, e models.Poll) (string, error)
	Delete(ctx context.Context, id, uID string) error
	//TODO add paginator to getMany
	GetMany(ctx context.Context, uID string) ([]models.Poll, error)
	Get(ctx context.Context, id string) (models.Poll, error)
}

type basicService struct {
	repo pollrepo.Repository
}

// NewService create a new service
func NewService(repo pollrepo.Repository) Service {
	return basicService{
		repo: repo,
	}
}

// Create create anew poll/poll
func (s basicService) Create(ctx context.Context, p models.Poll) (string, error) {
	// set default
	p = setupPoll(p)
	return p.ShortCode, s.repo.Create(ctx, p)
}

func (s basicService) Delete(ctx context.Context, id, uID string) error {
	return s.repo.Delete(ctx, id, uID)
	// TODO delete poll from cache
}

func (s basicService) GetMany(ctx context.Context, uID string) ([]models.Poll, error) {
	return s.repo.GetMany(ctx, uID)
}

func (s basicService) Get(ctx context.Context, id string) (models.Poll, error) {
	return s.repo.Get(ctx, id)
}

// sets up the poll
func setupPoll(p models.Poll) models.Poll {

	pll := p

	pll.ShortCode = str.WriteFromFormat("ddaAAAA")
	pll.CreatedAt = time.Now()

	// default identity is ip address
	if pll.Identity == "" {
		pll.Identity = models.IDIPAddr
	}

	// if poll is anonymous
	if pll.UserID == primitive.NilObjectID {
		et := time.Now().Add(time.Hour * 24)
		pll.ExpiresAt = &et
	}

	pll.Summary = models.PollSummary{
		OptionCount: make([]int64, len(pll.Option)),
	}

	return pll
}
