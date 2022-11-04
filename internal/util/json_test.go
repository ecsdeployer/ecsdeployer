package util

import (
	"testing"
)

type jsonStruct struct {
	Thing string `json:"thing"`
}

func TestJsonify(t *testing.T) {

	tables := []struct {
		obj interface{}
	}{
		{5},
		{"test"},
		{true},
		{jsonStruct{Thing: "blahblahba"}},
		{&jsonStruct{Thing: "blahblahba"}},
		{[]int{1, 2, 3, 4, 5}},
		{nil},
	}

	for _, table := range tables {
		res, err := Jsonify(table.obj)
		if err != nil {
			t.Fatalf("Failed to Jsonify: %s", err)
		}
		// ensure interface
		var _ string = res
	}
}

func TestJsonifyPretty(t *testing.T) {

	tables := []struct {
		obj interface{}
	}{
		{5},
		{"test"},
		{true},
		{jsonStruct{Thing: "blahblahba"}},
		{&jsonStruct{Thing: "blahblahba"}},
		{[]int{1, 2, 3, 4, 5}},
		{nil},
	}

	for _, table := range tables {
		res, err := JsonifyPretty(table.obj)
		if err != nil {
			t.Fatalf("Failed to JsonifyPretty: %s", err)
		}
		// ensure interface
		var _ string = res
	}
}
