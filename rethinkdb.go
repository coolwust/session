package session

import (
	r "github.com/dancannon/gorethink"
	"time"
)

var _ Store = new(RethinkDBStore)

type RethinkDBStore struct {
	Session *r.Session
	DB      string
	Table   string
}

func NewRethinkDBStore(dbSession *r.Session, db, table string) *RethinkDBStore {
	return &RethinkDBStore{
		Session: dbSession,
		DB:      db,
		Table:   table,
	}
}

func (store *RethinkDBStore) Set(session *Session) error {
	_, err := r.DB(store.DB).Table(store.Table).Insert(map[string]interface{}{
		"id":       session.ID,
		"expires":  session.Expires.UnixNano(),
		"duration": session.Duration,
		"data":     session.Data,
	}).RunWrite(store.Session)
	return err
}

func (store *RethinkDBStore) Get(id string) (*Session, error) {
	cur, err := r.DB(store.DB).Table(store.Table).Get(id).Run(store.Session)
	if err != nil {
		return nil, err
	}
	if cur.IsNil() {
		return nil, ErrNotFound
	}
	m := make(map[string]interface{})
	if err := cur.One(&m); err != nil {
		return nil, err
	}
	return &Session{
		ID:       m["id"].(string),
		Expires:  time.Unix(0, int64(m["expires"].(float64))),
		Duration: time.Duration(m["duration"].(float64)),
		Data:     m["data"].(map[string]interface{}),
	}, nil
}

func (store *RethinkDBStore) Remove(id string) error {
	_, err := r.DB(store.DB).Table(store.Table).Get(id).Delete().RunWrite(store.Session)
	return err
}

func (store *RethinkDBStore) Clean() error {
	table := r.DB(store.DB).Table(store.Table)
	now := time.Now().UnixNano()
	_, err := table.Filter(r.Row.Field("expires").Le(now)).Delete().RunWrite(store.Session)
	return err
}
