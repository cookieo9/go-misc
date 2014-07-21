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
	vch := reflect.ValueOf(ch)

	if vch.Kind() != reflect.Chan { // Not a channel
		return errNotChan
	}

	if vch.Type().ChanDir()&reflect.RecvDir == 0 { // Can't Receive
		return errNotRecv
	}

	p.input = vch
	p.subs = make(chan sub)

	return nil
}

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
			vch := reflect.ValueOf(ch)

			if vch.Kind() != reflect.Chan { // Not a channel
				subscription.err <- errNotChan
				break
			}

			if vch.Type().ChanDir()&reflect.SendDir == 0 { // Can't send
				subscription.err <- errNotSend
				break
			}

			if !p.input.Type().Elem().AssignableTo(vch.Type().Elem()) { // Bad channel type
				subscription.err <- errBadType
				break
			}

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
