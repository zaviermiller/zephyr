package zephyr

import (
	"fmt"
	"runtime"
	"strings"
)

// Interfaces =-=
// The following are interfaces that are implemented by
// the reactive data.

// Listener lets implementations call an Update()
// function, which triggers a DOM update in the
// reactive data types
type Listener interface {
	Update()
	Identifier() string
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
	Listeners map[string]Listener
}

// void is for internal use in the simple set
// implementation
var void struct{}

// NewRD creates and returns a new ReactiveData
// struct with the given data.
func NewRD(data interface{}) ReactiveData {
	var rd ReactiveData
	rd = ReactiveData{Data: data, Listeners: map[string]Listener{}}

	return rd
}

// RegisterOnComponent registers a component-wide listener that will
// trigger a whole vDOM re-render/update
// func (rd *ReactiveData) RegisterOnComponent(l *ComponentListener) {
// 	rd.Listeners[l.ID] = Listener(l)
// }

func (rd *ReactiveData) Notify() {
	for _, l := range rd.Listeners {
		l.Update()
	}
}

// RegisterOnNode

// Register is handles and registers the various
// listener implementations.
func (rd *ReactiveData) Register(l Listener) {
	if l != nil {
		rd.Listeners[l.Identifier()] = l
	}
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
	Value(l Listener) interface{}

	// String is used by the vDOM to render HTML
	// easily. All types should have this, which
	// allows for clean use in the Render() func.
	string(l Listener) string
}

// ZephyrString is the ZephyrData implementation
// for the `string` type.
type ZephyrString func() *ReactiveData

// NewLiveString returns a "live" string (reactive type ZephyrString)
// Change to NewZephyrString?
func (c *BaseComponent) NewLiveString(data string) ZephyrString {
	// create a new ReactiveData
	rd := NewRD(data)
	rdPtr := &rd
	fmt.Println(rd.Data)
	// return func type with getter
	rdGetter := ZephyrString(func() *ReactiveData {
		return rdPtr
	})
	return rdGetter
}

// Set implements ZephyrData.Set(interface{}),
// and is used to set and notify listeners.
func (str ZephyrString) Set(newData interface{}) {
	val, ok := newData.(string)
	if !ok {
		panic("invalid data type - fixme")
	}
	// setter func?
	rd := str()
	rd.Data = val
	fmt.Println("set: ", rd.Data)
	fmt.Println("notifying children: ", rd.Listeners)
	rd.Notify()
}

func (str ZephyrString) Value(l Listener) interface{} {
	fmt.Println(l)
	pc, _, _, _ := runtime.Caller(1)
	funcName := runtime.FuncForPC(pc).Name()
	lastSlash := strings.LastIndexByte(funcName, '/')
	if lastSlash < 0 {
		lastSlash = 0
	}
	lastDot := strings.LastIndexByte(funcName[lastSlash:], '.') + lastSlash

	fmt.Printf("Package: %s\n", funcName[:lastDot])
	fmt.Printf("Func:   %s\n", funcName[lastDot+1:])
	rd := str()
	rd.Register(l)

	return str().Data
}

// String implements Zephyr.String() string,
// and is used internally by the HTML renderer
func (str ZephyrString) string(l Listener) string {
	rd := str()
	rd.Register(l)
	fmt.Println("got string, set listener ", rd.Listeners, l)
	return rd.Data.(string)
}

// todo
type ComponentListener struct {
	Updater func()
	id      string
}

func (l ComponentListener) Update() {
	// re-render component on update
	if l.Updater != nil {
		l.Updater()
	}
}

func (l ComponentListener) Identifier() string {
	return l.id
}

// func ()

type VNodeListener struct {
	Updater func()
	id      string
}

func (l VNodeListener) Update() {
	// re-render component on update
	if l.Updater != nil {
		l.Updater()
	}
}

func (l VNodeListener) Identifier() string {
	return l.id
}
