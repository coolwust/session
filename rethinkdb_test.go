package session

import (
	"testing"
	"log"
	"os"
	"time"
	"reflect"
	r "github.com/dancannon/gorethink"
)

var (
	rSession *r.Session
	rAddress string
	rDB      string
	rTable   string
)


func init() {
	rAddress = os.Getenv("SESSION_RETHINKDBSTORE_ADDRESS")
	rDB      = os.Getenv("SESSION_RETHINKDBSTORE_DB")
	rTable   = os.Getenv("SESSION_RETHINKDBSTORE_TABLE")

	var err error
	rSession, err = r.Connect(r.ConnectOpts{
		Address:  rAddress,
		Database: rDB,
	})
	if err != nil {
		log.Fatalln(err)
	}
}

func rSetUp() {
	r.Table(rTable).Delete().RunWrite(rSession)
}

func rTearDown() {
	r.Table(rTable).Delete().RunWrite(rSession)
}

func TestRethinkDBStore(t *testing.T) {
	rSetUp()
	defer rTearDown()
	store := NewRethinkDBStore(rSession, rDB, rTable)
	session1 := NewSession(UUID, time.Hour)
	session1.Set("foo", "bar")
	if err := store.Set(session1); err != nil {
		t.Fatal(err)
	}
	session2, err := store.Get(session1.ID)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(session1, session2) {
		t.Fatal("expected %v, got %v", session1, session2)
	}
	if err := store.Remove(session1.ID); err != nil {
		t.Fatal(err)
	}
	if _, err := store.Get(session1.ID); err != ErrNotFound {
		t.Fatal("expected the session is removed from storage")
	}
}

func TestRethinkDBStoreClean(t *testing.T) {
	rSetUp()
	defer rTearDown()
	store := NewRethinkDBStore(rSession, rDB, rTable)
	session1 := &Session{ID: "1", Expires: time.Now().Add(-time.Hour)}
	session2 := &Session{ID: "2", Expires: time.Now().Add(time.Hour)}
	store.Set(session1)
	store.Set(session2)
	if err := store.Clean(); err != nil {
		t.Fatal(err)
	}
	if _, err := store.Get(session1.ID); err != ErrNotFound {
		t.Fatal("session1 should be deleted")
	}
	if _, err := store.Get(session2.ID); err == ErrNotFound {
		t.Fatal("session2 should not be deleted")
	}
}
