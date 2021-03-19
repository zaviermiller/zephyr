package zephyr

import (
	"fmt"
	"reflect"
)

// Interfaces =-=
// The following are interfaces that are implemented by
// the reactive data.

// Listener lets implementations call an Update()
// function, which triggers a DOM update in the
// reactive data types
type Listener interface {
	Update()
}

// Subject lets implementations register
// listeners and notify their listeners.
// May want to add removal in the future.
type Subject interface {
	Register(l Listener)
	Notify()
}

// ZephyrData
// NewReactive__ -> Zephyr__ (ZephyrData impls) -- func() __
// Zephyr__:
//

// IMPLEMENTATIONS =-=-=

// ReactiveData is the struct that acts as
// the struct implementation of the Subject.
type ReactiveData struct {
	Data      interface{}
	Listeners map[Listener]struct{}
}

// void is for internal use in the simple set
// implementation
var void struct{}

// NewRD creates and returns a new ReactiveData
// struct with the given data.
func NewRD(data interface{}) ReactiveData {
	var rd ReactiveData
	rd = ReactiveData{Data: data, Listeners: map[Listener]struct{}{}}

	return rd
}

// RegisterOnComponent registers a component-wide listener that will
// trigger a whole vDOM re-render/update
// func (rd *ReactiveData) RegisterOnComponent(l *ComponentListener) {
// 	rd.Listeners[l.ID] = Listener(l)
// }

func (rd *ReactiveData) Notify() {
	for l, _ := range rd.Listeners {
		l.Update()
	}
}

// RegisterOnNode

// Register is handles and registers the various
// listener implementations.
func (rd *ReactiveData) Register(l Listener) {
	switch l.(type) {
	case VNodeListener:
		// rd.RegisterOnComponent(l.(*ComponentListener))
		rd.Listeners[l] = void
	default:
		fmt.Println(reflect.TypeOf(l))
	}
}

// ZephyrDepType is the type of dependency using the
// data. Internal is mostly for internal use, but can
// give more control over ReactiveData. WARNING: Internal
// deps are not reactive!!
type InternalListener struct{}

func (l InternalListener) Update() {

}

// ZephyrData is the interface for using
// reactive data in components. Implementations
// are functions that return their resp. types.
type ZephyrData interface {
	// Set must be implemented on all data types.
	// Type checking occurs in implementations.
	Set(interface{})

	// Value returns the value stored inside
	// the reactive data; requires type assert
	Value() interface{}

	// String is used by the vDOM to render HTML
	// easily. All types should have this, which
	// allows for clean use in the Render() func.
	String() string
}

// ZephyrString is the ZephyrData implementation
// for the `string` type.
type ZephyrString func(Listener) *ReactiveData

// NewLiveString returns a "live" string (reactive type ZephyrString)
// Change to NewZephyrString?
func (c *BaseComponent) NewLiveString(data string) ZephyrString {
	// create a new ReactiveData
	rd := NewRD(data)
	// return func type with getter
	strGetter := ZephyrString(func(l Listener) *ReactiveData {
		switch l.(type) {
		case ComponentListener:
			fmt.Println("node render listener here")
			rd.Register(l)
			return &rd
		case VNodeListener:
			rd.Register(l)
		case InternalListener:
			return &rd
		default:
			panic("context not recognized")
		}
		return &rd
	})
	return strGetter
}

// Set implements ZephyrData.Set(interface{}),
// and is used to set and notify listeners.
func (str ZephyrString) Set(newData interface{}) {
	val, ok := newData.(string)
	if !ok {
		panic("invalid data type - fixme")
	}
	// setter func?
	rd := str(InternalListener{})
	rd.Data = val
	rd.Notify()
}

func (str ZephyrString) Value() interface{} {
	return str(InternalListener{}).Data
}

// String implements Zephyr.String() string,
// and is used internally by the HTML renderer
func (str ZephyrString) String() string {
	val := str(InternalListener{})
	return val.Data.(string)
}

// todo
type ComponentListener struct {
	Updater func()
}

func (l ComponentListener) Update() {
	// re-render component on update
	if l.Updater != nil {
		l.Updater()
	}
}

// func ()

type VNodeListener struct {
	Updater func()
}

func (l VNodeListener) Update() {
	// re-render component on update
	if l.Updater != nil {
		l.Updater()
	}
}
