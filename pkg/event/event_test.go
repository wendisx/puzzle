package event

import (
	"strconv"
	"testing"

	"github.com/wendisx/puzzle/pkg/config"
)

// [passed]
func Test_obserable_event(t *testing.T) {
	config.LoadDict(config.DICTKEY_EVENT)
	eventDict := config.GetDict(config.DICTKEY_EVENT)
	listener := &ObserableListener{
		id: 1,
	}
	eventDict.Record(feature_listener+strconv.Itoa(int(listener.id)), listener)
	ep := ObserablePublisher{}
	e := NewObserableEvent(999, nil)
	ep.Publish(e)
}

func Test_unobserable_event(t *testing.T) {
	config.LoadDict(config.DICTKEY_EVENT)
	eventDict := config.GetDict(config.DICTKEY_EVENT)
	listener := &ObserableListener{
		id: 2,
	}
	eventDict.Record(feature_listener+strconv.Itoa(int(listener.id)), listener)
	ep := ObserablePublisher{}
	e := NewUnobserableEvent()
	ep.Publish(e)
}
