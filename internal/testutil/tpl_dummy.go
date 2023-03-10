package testutil

import "errors"

// this fakes the interface to a Template object
type tplDummy struct{}

func (*tplDummy) Apply(val string) (string, error) {
	return val, nil
}

var TplDummy = &tplDummy{}

// failure dummy
var ErrTplDummyFailureError = errors.New("tpl dummy failure")

type tplDummyFailure struct{}

func (*tplDummyFailure) Apply(val string) (string, error) {
	return "", ErrTplDummyFailureError
}

var TplDummyFailure = &tplDummyFailure{}
