package auth

import (
	"errors"
	"net/http"
	"reflect"
	"testing"
)

func TestNoAuth(t *testing.T) {
	header := http.Header{}
	_, got := GetAPIKey(header)
	want := errors.New("no authorization header included")
	if !reflect.DeepEqual(want, got) {
		t.Fatalf("expected: %v, got: %v", want, got)
	}
}

func TestAuthFormat(t *testing.T) {
	header := http.Header{
		"Authorization": {"login"},
	}
	_, got := GetAPIKey(header)
	want := errors.New("malformed authorization header")
	if !reflect.DeepEqual(want, got) {
		t.Fatalf("expected: %v, got: %v", want, got)
	}
}

func TestSuccess(t *testing.T) {
	header := http.Header{
		"Authorization": {"ApiKey @jnnasdo1.,3214naJ!@HBFJS@IJIU$#@NUIANDSfn23"},
	}
	got, _ := GetAPIKey(header)
	want := "@jnnasdo1.,3214naJ!@HBFJS@IJIU$#@NUIANDSfn23"
	if !reflect.DeepEqual(want, got) {
		t.Fatalf("expected: %v, got: %v", want, got)
	}
}
