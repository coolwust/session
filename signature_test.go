package session

import (
	"testing"
)

var signTests = []struct {
	unsigned string
	key      string
	signed   string
}{
	{"hello", "foo", "hello.NAJWxCdrLTSD8CvzyLdIDhLH6pcCsiAKldMySgs4"},
	{"world", "bar", "world.6kpaLfNqzist38XcIeU9ejULKGCJK7u23K9Qwbyb4sk"},
}

func TestSign(t *testing.T) {
	for i, test := range signTests {
		if signed := Sign(test.unsigned, test.key); signed != test.signed {
			t.Errorf("%d:\n  %s + %s = %s (got %s)", i, test.unsigned, test.key, test.signed, signed)
		}
	}
}

var unsignTests = []struct {
	signed   string
	key   string
	unsigned string
}{
	{"hello.NAJWxCdrLTSD8CvzyLdIDhLH6pcCsiAKldMySgs4", "foo", "hello"},
	{"world.6kpaLfNqzist38XcIeU9ejULKGCJK7u23K9Qwbyb4sk", "bar", "world"},
	{"hello", "foo", ""},
	{"world.xxx", "bar", ""},
}

func TestUnsign(t *testing.T) {
	for i, test := range unsignTests {
		if unsigned, _ := Unsign(test.signed, test.key); unsigned != test.unsigned {
			t.Errorf("%d:\n %s + %s = %s (got %s)", i, test.signed, test.key, test.unsigned, unsigned)
		}
	}
}
