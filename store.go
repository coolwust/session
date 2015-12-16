package session

import (
	"errors"
)

var ErrNotFound = errors.New("The session is not found in storage")

type Store interface {

	Set(session *Session) error

	Get(id string) (*Session, error)

	Remove(id string) error

	Clean() error
}
