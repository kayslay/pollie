package voteendpoint

import "errors"

type Vote struct {
	Option    []int  `json:"option" bson:"option"`
	IP        string `json:"-" bson:"ip"`
	PollID    string `json:"-" bson:"poll_id"`
	Device    string `json:"-" bson:"device"`
	UserAgent string `json:"-" bson:"user_agent"`
}

func (v *Vote) Validate(id, ip string) error {

	v.IP, v.PollID = ip, id

	if len(v.PollID) != 7 {
		return errors.New("poll code is of length 7")
	}

	return nil
}
