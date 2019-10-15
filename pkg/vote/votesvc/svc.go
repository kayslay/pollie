package votesvc

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"pollie"
	"pollie/models"
	"pollie/pkg/poll/pollrepo"
	"pollie/pkg/vote/voterepo"
	"time"

	"github.com/go-redis/redis"
	"github.com/pkg/errors"
)

var (
	errOptionRange = errors.New("option_range: option passed is invalid")
)

type Service interface {
	Vote(ctx context.Context, v models.Vote) (models.Poll, error)
	GetPoll(ctx context.Context, id string) (models.Poll, error)
}

func NewService(v voterepo.Repository, p pollrepo.Repository, r *redis.Client) Service {
	return basicService{repo: v, pollRepo: p, redis: r}
}

type basicService struct {
	repo     voterepo.Repository
	pollRepo pollrepo.Repository
	redis    *redis.Client
}

func (s basicService) Vote(ctx context.Context, v models.Vote) (models.Poll, error) {
	// get poll
	p := ctx.Value(pollie.ContextKey("poll")).(models.Poll)

	// check if the user can vote
	if err := s.validateIdentity(p.ID.Hex(), p.Identity, v.Meta.IP); err != nil {
		return p, err
	}

	// validate the vote
	err := s.validateVote(p, &v)
	if err != nil {
		return p, err
	}
	// set the meta of the vote
	setMeta(&v.Meta)

	// create the vote
	_, err = s.repo.Vote(v)
	if err != nil {
		return p, pollie.WrapErr(err, "create_vote")
	}

	// update the poll summary
	log.Println("update poll summary", s.repo.UpdatePollSummary(p.ID.Hex(), v.Option))

	// get the updated poll document
	_p, err := s.pollRepo.Get(ctx, p.ID.Hex())
	if err != nil {
		// return previous poll
		return p, nil
	}

	return _p, err
}

func (s basicService) GetPoll(ctx context.Context, id string) (models.Poll, error) {

	var poll models.Poll
	redisKey := fmt.Sprintf("poll:%s", id)

	// check the cache for poll
	sPoll, err := s.redis.Get(redisKey).Result()
	if err == nil {

		err = json.Unmarshal([]byte(sPoll), &poll)
		if err == nil {
			return poll, nil
		}
	}

	//
	poll, err = s.pollRepo.Get(ctx, id)
	// error getting the poll
	if err != nil {
		return poll, err
	}

	b, err := json.Marshal(poll)
	expire := time.Duration(time.Minute * 30)
	// use the lowest time
	if poll.ExpiresAt != nil {

		timeLeft := poll.ExpiresAt.Sub(time.Now())
		// if the timeLeft is less that expire duration use
		// timeLeft to expire the cache

		if timeLeft < expire {
			if timeLeft < 0 {
				timeLeft = 1
			}
			expire = timeLeft
		}

		fmt.Println("---", timeLeft, expire)

	}
	if err == nil {
		s.redis.Set(redisKey, string(b), expire)
	}

	return poll, nil
}

// validateIdentity validates if the user is allowed to vote.
// checks if the poll uses ip, cookies, user, secret, social media e.t.c
func (s basicService) validateIdentity(eID, eIdentity, ip string) error {

	// return nil //TODO uncomment to check ip

	switch eIdentity {
	case models.IDIPAddr:
		b, err := s.repo.IPExists(eID, ip)
		if err != nil {
			return pollie.WrapErr(err, "validate_ip_addr")
		}

		if b == true {
			return errors.New("ip address exist")
		}

		return nil
	}

	return errors.New("unknown poll identity")
}

func setMeta(meta *models.Meta) {
	ipInfo, err := pollie.ForeignIP(meta.IP)
	if err != nil {
		return
	}

	// set meta
	meta.CountryName = ipInfo.CountryName
	meta.CountryCode = ipInfo.CountryCode
	meta.RegionName = ipInfo.RegionName
	meta.RegionCode = ipInfo.RegionCode
	meta.Latitude = ipInfo.Latitude
	meta.Longitude = ipInfo.Longitude

}

func (s basicService) validateVote(p models.Poll, v *models.Vote) error {
	switch p.Type {
	case models.PTypeSingle:
		if len(v.Option) == 0 {
			if !p.NilVote {
				return errors.New("empty vote is note allowed")
			}
		} else {
			// pass only the first option
			v.Option = v.Option[0:1]
			if v.Option[0] >= len(p.Option) || v.Option[0] < 0 {
				return errOptionRange
			}

			return nil
		}

	default:
		return models.ErrInvalidPollType
	}

	return models.ErrInvalidPollType
}
