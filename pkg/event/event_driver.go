package event

import (
	"fmt"
	"strings"
	"time"

	"github.com/wendisx/puzzle/pkg/clog"
	"github.com/wendisx/puzzle/pkg/config"
)

/*
	event driver
	1. 接口: 事件, 事件发布者, 事件监听者
	2. 实现: 可观测事件, 可观测事件发布者, 可观测事件监听者
	3. 监听者构建后自动进入数据字典, 一般情况下不过期, 作为一种及时function call的存在执行一些必要的处理逻辑
	4. 数据字典中可能存在非监听者的数据key, 这种情况下需要考虑标识
*/

const (
	event_timeline    = "2006-01-02 15:04:05"
	feature_event     = "event:"
	feature_publisher = "publisher:"
	feature_listener  = "listener:"
)

const (
	EVENT_KIND_UNOBSERABLE = 0
	EVENT_KIND_OBSERABLE   = 1
)

type (
	Event interface {
		Kind() uint
		Timeline() string
		Rarity() int
		Record() any
	}
	EventPublisher interface {
		Publish(e Event) error
	}
	EventListener interface {
		Want(e Event) bool
		Listen(e Event) error
	}
	EventMeta struct {
		timeline time.Time
		rarity   int // for some specific operations
	}
	// maybe Use of some meta information or more weird states
	UnobserableEvent struct {
		EventMeta
	}
	ObserableEvent struct {
		EventMeta
		record any
	}
	ObserablePublisher struct {
	}
	ObserableListener struct {
		id uint
	}
)

func NewObserableEvent(rarity int, record any) *ObserableEvent {
	return &ObserableEvent{
		EventMeta: EventMeta{
			timeline: time.Now(),
			rarity:   rarity,
		},
		record: record,
	}
}

func NewUnobserableEvent() *UnobserableEvent {
	return &UnobserableEvent{
		EventMeta: EventMeta{
			timeline: time.Now(),
			rarity:   -1,
		},
	}
}

func (e ObserableEvent) Kind() uint {
	return EVENT_KIND_OBSERABLE
}

func (e ObserableEvent) Timeline() string {
	return e.timeline.Format(event_timeline)
}

func (e ObserableEvent) Rarity() int {
	return e.rarity
}

func (e ObserableEvent) Record() any {
	return e.record
}

func (e UnobserableEvent) Kind() uint {
	return EVENT_KIND_UNOBSERABLE
}

func (e UnobserableEvent) Timeline() string {
	return e.timeline.Format(event_timeline)
}

func (e UnobserableEvent) Rarity() int {
	return e.rarity
}

func (e UnobserableEvent) Record() any {
	return nil
}

func (p *ObserablePublisher) Publish(e Event) error {
	eventDict := config.GetDict(config.DICTKEY_EVENT)
	filter := func(k string) bool {
		return strings.HasPrefix(k, feature_listener)
	}
	keys := eventDict.Keys(filter)
	var listener EventListener
	for _, k := range keys {
		listener = eventDict.Find(k).Value().(EventListener)
		if listener.Want(e) {
			listener.Listen(e)
		}
	}
	return nil
}

func (l *ObserableListener) Want(e Event) bool {
	return e.Kind()&EVENT_KIND_OBSERABLE == EVENT_KIND_OBSERABLE
}

func (l *ObserableListener) Listen(e Event) error {
	clog.Info(fmt.Sprintf("[%d]event listen: (+%s, %d, T: %d)", l.id, e.Timeline(), e.Rarity(), e.Kind()))
	return nil
}
