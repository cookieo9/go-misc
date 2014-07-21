package chans_test

import (
	"fmt"
	"log"

	"github.com/cookieo9/go-misc/chans"
)

func ExamplePublisher() {
	in := make(chan int)
	a, b, c := make(chan int), make(chan int), make(chan int)

	pub, err := chans.NewPublisher(in)
	if err != nil {
		log.Fatal(err)
	}

	if err := pub.Subscribe(a); err != nil {
		log.Fatal(err)
	}

	if err := pub.Subscribe(b); err != nil {
		log.Fatal(err)
	}

	var x, y, z int

	in <- 5
	x, y = <-a, <-b
	fmt.Printf("x = %d ; y = %d ; z = %d\n", x, y, z)

	if err := pub.Subscribe(c); err != nil {
		log.Fatal(err)
	}
	if err := pub.Unsubscribe(b); err != nil {
		log.Fatal(err)
	}

	in <- 42
	in <- 6

	x, z = <-a, <-c
	fmt.Printf("x = %d ; y = %d ; z = %d\n", x, y, z)

	select {
	case y = <-b:
		log.Fatal("recieved a value on b when none expected!")
	default:
	}
	close(in)

	// Output:
	// x = 5 ; y = 5 ; z = 0
	// x = 6 ; y = 5 ; z = 6
}
