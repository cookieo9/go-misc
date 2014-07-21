// +build go1.1

package chans

import (
	"errors"
	"reflect"
	"sync"
)

var (
	errNotChan  = errors.New("not a channel")
	errNotRecv  = errors.New("can't receive from channel")
	errNotSend  = errors.New("can't send to channel")
	errBadType  = errors.New("bad channel type")
	errPubDead  = errors.New("publisher input channel closed")
	errSubExist = errors.New("already subscribed")
	errSubNone  = errors.New("not subscribed")
)

// A Publisher maintains a publish/subscribe relationship between
// an input channel, and a set of output channels.
type Publisher struct {
	input reflect.Value
	subs  chan sub
	dead  bool
	lock  sync.Mutex
}

type sub struct {
	unsub bool
	ch    interface{}
	err   chan error
}

// NewPublisher creates a new Publisher that reads messages from
// the given channel.
//
// Will return an error if the ch argument is not a channel or
// cannot be received from.
func NewPublisher(ch interface{}) (*Publisher, error) {
	pub := new(Publisher)

	if err := pub.init(ch); err != nil {
		return nil, err
	}

	pub.lock.Lock()
	go pub.main()

	return pub, nil
}

func (p *Publisher) init(ch interface{}) error {
	vch, err := p.checkChannel(ch, true)
	if err != nil {
		return err
	}

	p.input = vch
	p.subs = make(chan sub)

	return nil
}

// Unsubscribe stops the given channel from recieving messages
// sent through the Publisher.
//
// Returns an error if the channel was never subscribed to the
// Publisher.
func (p *Publisher) Unsubscribe(ch interface{}) error {
	p.lock.Lock()
	defer p.lock.Unlock()

	if p.dead {
		return errPubDead
	}

	ech := make(chan error)
	defer close(ech)

	p.subs <- sub{
		unsub: true,
		ch:    ch,
		err:   ech,
	}
	return <-ech
}

// Subscribe adds the given channel to the Publisher's list of subscribers
// and will begin to receive messages sent to the Publisher's input channel.
//
// Returns an error if ch:
//  - is not a channel
//  - cannot be sent to
//  - has a different element type than the input channel
func (p *Publisher) Subscribe(ch interface{}) error {
	p.lock.Lock()
	defer p.lock.Unlock()

	if p.dead {
		return errPubDead
	}

	ech := make(chan error)
	defer close(ech)

	p.subs <- sub{
		ch:  ch,
		err: ech,
	}
	return <-ech
}

func (p *Publisher) checkChannel(ch interface{}, input bool) (reflect.Value, error) {
	vch := reflect.ValueOf(ch)
	vnil := reflect.ValueOf(nil)

	if vch.Kind() != reflect.Chan {
		return vnil, errNotChan
	}

	if input {
		if vch.Type().ChanDir()&reflect.RecvDir == 0 {
			return vnil, errNotRecv
		}
		return vch, nil
	}

	if vch.Type().ChanDir()&reflect.SendDir == 0 {
		return vnil, errNotSend
	}

	if !p.input.Type().Elem().AssignableTo(vch.Type().Elem()) {
		return vnil, errBadType
	}

	return vch, nil
}

func (p *Publisher) main() {
	var (
		subscribers   []reflect.Value
		subscriberSet = map[interface{}]bool{}
		waiting       []reflect.SelectCase
	)
	p.lock.Unlock()

	inputCase := reflect.SelectCase{
		Dir:  reflect.SelectRecv,
		Chan: p.input,
	}
	subCase := reflect.SelectCase{
		Dir:  reflect.SelectRecv,
		Chan: reflect.ValueOf(p.subs),
	}

	waiting = append(waiting, inputCase)
	waiting = append(waiting, subCase)

	for {
		idx, val, ok := reflect.Select(waiting)

		switch idx {
		case 0: // New message
			if !ok { // input channel closed
				p.lock.Lock()
				defer p.lock.Unlock()

				close(p.subs)
				p.dead = true

				return
			}

			// Send message to all current subscribers
			waiting = waiting[:2]
			for _, subscriber := range subscribers {
				waiting = append(waiting, reflect.SelectCase{
					Dir:  reflect.SelectSend,
					Chan: subscriber,
					Send: val,
				})
			}

		case 1: // New subscriber
			subscription := val.Interface().(sub)
			ch := subscription.ch

			if subscription.unsub { // Removing a subscription
				if !subscriberSet[ch] { // No match
					subscription.err <- errSubNone
					break
				}

				// Find and remove
				for idx, subscriber := range subscribers {
					if subscriber.Interface() == ch {
						delete(subscriberSet, ch)

						n := len(subscribers)
						subscribers[idx] = subscribers[n-1]
						subscribers = subscribers[:n-1]
						break
					}
				}
			} else {
				vch, err := p.checkChannel(ch, false)
				if err != nil { // bad channel
					subscription.err <- err
					break
				}

				if subscriberSet[ch] { // Already exists
					subscription.err <- errSubExist
					break
				} else { // Add subscription
					subscriberSet[ch] = true
					subscribers = append(subscribers, vch)
				}
			}
			subscription.err <- nil

		default: // Value sent
			n := len(waiting)
			waiting[idx] = waiting[n-1]
			waiting = waiting[:n-1]
		}
	}
}
