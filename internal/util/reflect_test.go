package util

import (
	"testing"
)

type fooS struct {
	SomeVal int32
	StringP *string
	StringV string
}

type barS struct {
	NestA *barS
	NestB *fooS
	Foo   *fooS
	Bar   *barS
	Blah1 *string
	Blah2 *int32
	Blah3 int32
}

type nester struct {
	Nest1   *nester
	Nest2   *nester
	Foo     *fooS
	Bar     *barS
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
					Foo: &fooS{
						StringV: "test",
					},
				},
				Nest2: &nester{
					Foo: &fooS{
						StringV: "test2",
					},
					Bar: &barS{
						NestB: &fooS{
							StringV: "test4",
						},
					},
				},
				Foo: &fooS{
					StringV: "test3",
				},
			},
		},
	}

	_ = DeepFindInStruct[fooS](commonHaystack)
	// if len(stuff) != 4 {
	// 	t.Errorf("Not enough results! got=%d", len(stuff))
	// }
	// fmt.Printf("STUFF FOUND: %v\n", stuff)

	// for _, entry := range stuff {
	// 	// fmt.Printf("StringV = %v\n", entry.StringV)
	// 	_ = entry
	// }

	// tables := []struct {
	// 	haystack   *parentThing
	// 	finderFunc func(interface{}) []interface{}
	// }{
	// 	{commonHaystack, func(thing interface{}) []any { return DeepFindInStruct[fooS](thing).([]any)} },
	// }

	// for _, table := range tables {
	// 	result := table.finderFunc(commonHaystack)
	// 	if result != table.serviceName {
	// 		t.Errorf("expected <%s> to give service name of <%s> but got <%s>", table.arn, table.serviceName, result)
	// 	}
	// }
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
					Foo: &fooS{
						StringV: "test",
					},
				},
			},
		},
	}

	stuff := DeepFindInStruct[fooS](commonHaystack)
	if len(stuff) != 5 {
		t.Errorf("Not enough results! got=%d", len(stuff))
	}
	// fmt.Printf("STUFF FOUND: %v\n", stuff)

	for _, entry := range stuff {
		// fmt.Printf("StringV = %v\n", entry.StringV)
		_ = entry
	}
}
