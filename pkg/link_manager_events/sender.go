package link_manager_events

import (
	"log"

	om "github.com/adairxie/delinkcious/pkg/object_model"
	"github.com/nats-io/nats.go"
)

type eventSender struct {
	hostname string
	nats     *nats.EncodedConn
}

func (s *eventSender) OnLinkAdded(username string, link *om.Link) {
	err := s.nats.Publish(subject, Event{om.LinkAdded, username, link})
	if err != nil {
		log.Fatal(err)
	}
}

func (s *eventSender) OnLinkUpdated(username string, link *om.Link) {
	err := s.nats.Publish(subject, Event{om.LinkUpdated, username, link})
	if err != nil {
		log.Fatal(err)
	}
}

func (s *eventSender) OnLinkDeleted(username string, url string) {
	// Ignore link delete events
}

func NewEventSender(url string) (om.LinkManagerEvents, error) {
	ec, err := connect(url)
	if err != nil {
		return nil, err
	}
	return &eventSender{hostname: url, nats: ec}, nil
}
