package testutil

// this fakes the interface to a Template object
type tplDummy struct{}

func (*tplDummy) Apply(val string) (string, error) {
	return val, nil
}

var TplDummy = &tplDummy{}
