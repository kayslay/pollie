package pollie

import "github.com/pkg/errors"

type LoggableErr interface {
	Error() string
	Log() string
}

type err struct {
	err error
	msg string
}

func (e err) Error() string {
	return e.msg
}

func (e err) Log() string {
	return e.err.Error()
}

// WrapErr wraps an error with message and retruns Error interface
func WrapErr(e error, msg string) LoggableErr {
	if e == nil {
		return nil
	}

	return err{err: errors.Wrap(e, msg), msg: e.Error()}
}
