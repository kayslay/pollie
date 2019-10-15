package models

type MgoCloser interface {
	Close() error
}
