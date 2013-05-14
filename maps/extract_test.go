package maps

import (
	// "github.com/cookieo9/go-misc/slice"
	"../slice"
	"reflect"
	"testing"
)

func TestGetKeys(t *testing.T) {
	m := map[int]string{1: "one", 2: "two", 42: "life"}
	var s []int
	e := []int{1, 2, 42}

	GetKeys(m, &s)
	slice.Sort(s, func(a, b int) bool { return a <= b })
	t.Log(s)

	if len(s) != len(m) {
		t.Error("Bad Length for output slice", m, s)
	}

	if !reflect.DeepEqual(e, s) {
		t.Errorf("Expected %v, got %v", e, s)
	}
}

func TestGetVals(t *testing.T) {
	m := map[int]string{1: "one", 2: "two", 42: "life"}
	var s []string
	e := []string{"life", "one", "two"}

	GetVals(m, &s)
	slice.Sort(s, func(a, b string) bool { return a <= b })
	t.Log(s)

	if len(s) != len(m) {
		t.Error("Bad Length for output slice", m, s)
	}

	if !reflect.DeepEqual(e, s) {
		t.Errorf("Expected %v, got %v", e, s)
	}
}

func TestGetPairs(t *testing.T) {
	m := map[int]string{1: "one", 2: "two", 42: "life"}
	type pair struct {
		Key int
		Val string
	}
	var s []pair
	e := []pair{{1, "one"}, {2, "two"}, {42, "life"}}

	GetPairs(m, &s)
	slice.Sort(s, func(a, b pair) bool { return a.Key <= b.Key })
	t.Log(s)

	if len(s) != len(m) {
		t.Error("Bad Length for output slice", m, s)
	}

	if !reflect.DeepEqual(e, s) {
		t.Errorf("Expected %v, got %v", e, s)
	}
}
