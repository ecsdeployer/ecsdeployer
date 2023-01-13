package util

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type fooS struct {
	SomeVal int32
	StringP *string
	StringV string
}

type barS struct {
	NestA     *barS
	NestB     *fooS
	Foo       *fooS
	Bar       *barS
	Blah1     *string
	Blah2     *int32
	Blah3     int32
	something *int32
}

type nester struct {
	Nest1   *nester
	Nest2   *nester
	Foo2    *fooS
	Bar2    *barS
	BarV    barS
	ArrBar  []barS
	PArrBar []*barS
}

type parentThing struct {
	Thing1 nester
	Thing2 *nester
}

func TestDeepFindInStruct(t *testing.T) {

	commonHaystack := &parentThing{
		Thing1: nester{
			Nest1: &nester{
				Nest1: &nester{
					Foo2: &fooS{
						StringV: "test",
					},
				},
				Nest2: &nester{
					Foo2: &fooS{
						StringV: "test2",
					},
					Bar2: &barS{
						NestB: &fooS{
							StringV: "test4",
						},
					},
				},
				Foo2: &fooS{
					StringV: "test3",
				},
			},
		},
	}

	res := DeepFindInStruct[fooS](commonHaystack)
	require.Len(t, res, 4)
}

func TestDeepFindInStruct_WithArrays(t *testing.T) {

	commonHaystack := &parentThing{
		Thing1: nester{
			ArrBar: []barS{
				{
					Foo: &fooS{
						StringV: "thing",
					},
				},
				{
					Foo: &fooS{
						StringV: "thing2",
					},
				},
				{
					Bar: &barS{
						Blah3: 1234,
					},
				},
			},
			PArrBar: []*barS{
				{
					Foo: &fooS{
						StringV: "thingx",
					},
				},
				{
					Foo: &fooS{
						StringV: "thing2x",
					},
				},
				{
					Bar: &barS{
						Blah3: 1234,
					},
				},
			},
			Nest1: &nester{
				Nest1: &nester{
					Foo2: &fooS{
						StringV: "test",
					},
				},
			},
		},
	}

	require.Len(t, DeepFindInStruct[fooS](commonHaystack), 5)
	require.Len(t, DeepFindInStruct[barS](commonHaystack), 6)

	// it's 1 because the top level is a nester, so children are ignored
	require.Len(t, DeepFindInStruct[nester](commonHaystack), 1)
}
