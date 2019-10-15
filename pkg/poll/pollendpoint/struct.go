package pollendpoint

import (
	"errors"

	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
)

type CreateReq struct {
	Description string   `json:"description" bson:"description"`
	Options     []string `json:"options" bson:"options"`
	Tag         []string `json:"tag" bson:"tag"`
	// TODO add other fields as more features are added
	Auth bool `json:"auth" bson:"auth"`
}

func (c *CreateReq) Validate() error {
	vErr := validate.Validate(
		&validators.StringIsPresent{Name: "description", Field: c.Description},
		&validators.IntIsLessThan{
			Name:     "options",
			Field:    len(c.Options),
			Compared: 6,
			Message:  "number of option should be less than 7",
		},

		&validators.IntIsGreaterThan{
			Name:     "options",
			Field:    len(c.Options),
			Compared: 1,
			Message:  "number of options should be more than 1",
		},
		// TODO add more validations
	)

	if vErr.HasAny() {
		return errors.New(vErr.Error())
	}

	return nil
}

// DeleteReq delete request struct
type DeleteReq struct {
	ID     string
	UserID string
}

// GetManyReq get many request struct
type GetManyReq struct {
	UserID string
	Q      string
	Page   int
	Limit  int
}

// GetReq get request struct
type GetReq struct {
	ID     string //id or short_code
	UserID string
}
