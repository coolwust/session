package session

type Store interface {
	Set(session *Session) error
	Get(id string) (*Session, error)
	Remove(id string) error
	Clean() error
}
