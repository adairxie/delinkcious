package link_manager_events

import (
	om "github.com/adairxie/delinkcious/pkg/object_model"
)

type Event struct {
	EventType om.EventTypeEnum
	Username  string
	Link      *om.Link
}
