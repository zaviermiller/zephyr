package zephyr

import (
	"fmt"
	"reflect"
)

type Listener interface {
	Update()
}

type Subject interface {
	Register(l Listener)
	Notify()
}

// IMPLEMENTATIONS =-=-=

// ReactiveData is the struct that holds the data
// being listened to, and notifies listeners of
// any changes.
type ReactiveData struct {
	Data      interface{}
	Listeners map[string]Listener
}

func NewRD(data interface{}) ReactiveData {
	var rd ReactiveData
	rd = ReactiveData{Data: data, Listeners: map[string]Listener{}}

	return rd
}

func (rd *ReactiveData) RegisterOnComponent(l *ComponentListener) {
	rd.Listeners[l.ID] = Listener(l)
}

func (rd *ReactiveData) Register(l interface{}) {
	switch l.(type) {
	case *ComponentListener:
		rd.RegisterOnComponent(l.(*ComponentListener))
	default:
		fmt.Println(reflect.TypeOf(l))
	}
}

func (rd *ReactiveData) Notify() {
	for _, l := range rd.Listeners {
		l.Update()
	}
}

type ComponentListener struct {
	ID      string
	Updater func()
}

func (l ComponentListener) Update() {
	// re-render component on update
	if l.Updater != nil {
		l.Updater()
	}
}

// func ()
