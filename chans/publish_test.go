package chans

import (
	"reflect"
	"testing"
)

func newPublisher(x interface{}, t *testing.T) *Publisher {
	p, err := NewPublisher(x)
	if err != nil {
		t.Fatalf("NewPublisher(%#v): %v", x, err)
	}
	return p
}

var publisherInputTests = []struct {
	inputs []interface{}
	err    error
}{
	{[]interface{}{5, 6, "hello", 4.2, false}, errNotChan},
	{[]interface{}{make(chan<- int), make(chan<- bool)}, errNotRecv},
	{[]interface{}{make(<-chan int), make(<-chan bool)}, nil},
}

func TestPublisherInput(t *testing.T) {
	for _, test := range publisherInputTests {
		for _, input := range test.inputs {
			_, err := NewPublisher(input)
			if err == nil {
				defer reflect.ValueOf(input).Close()
			}
			if err != test.err {
				t.Errorf("NewPublisher(%v): expected %q got %q", input, test.err, err)
			} else {
				t.Logf("NewPublisher(%v): got expected error %q", input, err)
			}
		}
	}
}

type tFoo int

var publisherSubscribeTests = []struct {
	input, output interface{}
	err           error
}{
	{make(chan int), 5, errNotChan},
	{make(chan int), true, errNotChan},
	{make(chan int), make(chan bool), errBadType},
	{make(chan bool), make(chan bool), nil},
	{make(chan bool), make(<-chan bool), errNotSend},
	{make(chan tFoo), make(chan int), errBadType},
	{make(chan int), make(chan tFoo), errBadType},
	{make(chan tFoo), make(chan tFoo), nil},
}

func TestPublisherSubscribe(t *testing.T) {
	for _, test := range publisherSubscribeTests {
		p := newPublisher(test.input, t)
		defer reflect.ValueOf(test.input).Close()

		t.Logf("p(%T).Subscribe(%T) -> %#v", test.input, test.output, test.err)
		if err := p.Subscribe(test.output); err != test.err {
			t.Errorf("Got unexpected error: %#v", err)
		}
	}
}

func TestPublisherDead(t *testing.T) {
	a := make(chan int)
	b := make(chan int)

	p := newPublisher(a, t)
	close(a)

	if err := p.Subscribe(b); err != errPubDead {
		t.Fatalf("p.Subscribe(): expected %q, got %q", errPubDead, err)
	}
	if err := p.Unsubscribe(b); err != errPubDead {
		t.Fatalf("p.Unsubscribe(): expected %q, got %q", errPubDead, err)
	}
}

func TestPublisherSingleSubscriber(t *testing.T) {
	a := make(chan int)
	b := make(chan int)
	p := newPublisher(a, t)

	for _ = range make([]struct{}, 10) {
		if err := p.Subscribe(b); err != nil {
			t.Fatalf("p(%T).Subscribe(%T) -> %v", a, b, err)
		}

		if err := p.Subscribe(b); err != errSubExist {
			t.Fatalf("p(%T).Subscribe(%T) -> %v (expected: %v)", a, b, err, errSubExist)
		}

		if err := p.Unsubscribe(b); err != nil {
			t.Fatalf("p(%T).Unsubscribe(%T) -> %v", a, b, err)
		}

		if err := p.Unsubscribe(b); err != errSubNone {
			t.Fatalf("p(%T).Unsubscribe(%T) -> %v (expected: %v)", a, b, err, errSubNone)
		}
	}
}
