package main

import (
	"net/http"
	"testing"

	"github.com/anboo/codeception"
)

func TestSimple(t *testing.T) {
	a := codeception.NewActor(t, "http://beon.fun", make(map[string]string))
	a.SendGet("/", nil).SeeResponseCodeIs(http.StatusForbidden)
}
