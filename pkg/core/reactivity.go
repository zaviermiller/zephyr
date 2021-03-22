package zephyr

import (
	"fmt"
	"strconv"
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

// LiveData
// NewReactive__ -> Zephyr__ (LiveData impls) -- func() __
// Zephyr__:
//

// IMPLEMENTATIONS =-=-=

// DataDep is the struct that acts as
// the struct implementation of the Subject.
type DataDep struct {
	Data      interface{}
	Listeners map[string]Listener
}

// void is for internal use in the simple set
// implementation
var void struct{}

// NewDep creates and returns a new DataDep
// struct with the given data.
func NewDep(data interface{}) DataDep {
	var rd DataDep
	rd = DataDep{Data: data, Listeners: map[string]Listener{}}

	return rd
}

// RegisterOnComponent registers a component-wide listener that will
// trigger a whole vDOM re-render/update
// func (rd *DataDep) RegisterOnComponent(l *ComponentListener) {
// 	rd.Listeners[l.ID] = Listener(l)
// }

func (rd *DataDep) Notify() {
	for _, l := range rd.Listeners {
		l.Update()
	}
}

// RegisterOnNode

// Register is handles and registers the various
// listener implementations.
func (rd *DataDep) Register(l Listener) {
	if l != nil {
		rd.Listeners[l.Identifier()] = l
	}
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

// LiveData is the interface for using
// reactive data in components. Implementations
// are functions that return their resp. types.
type LiveData interface {
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

// Below are the implementations for all supported
// live data types.

// LiveString is the LiveData implementation
// for the `string` type.
type LiveString func() *DataDep

// NewLiveString returns a "live" string (reactive type LiveString)
func (c *BaseComponent) NewLiveString(data string) LiveString {
	// create a new DataDep
	rd := NewDep(data)
	rdPtr := &rd
	// return func type with getter
	rdGetter := LiveString(func() *DataDep {
		return rdPtr
	})
	return rdGetter
}

// Set implements LiveData.Set(interface{}),
// and is used to set and notify listeners.
func (str LiveString) Set(newData interface{}) {
	val, ok := newData.(string)
	if !ok {
		panic("invalid data type - fixme")
	}
	// setter func?
	rd := str()
	rd.Data = val
	rd.Notify()
}

func (str LiveString) Value(l Listener) interface{} {
	rd := str()
	if l == nil {
		// pc := make([]uintptr, 15)
		// n := runtime.Callers(2, pc)
		// frames := runtime.CallersFrames(pc[:n])
		// frame, _ := frames.Next()
		// fmt.Printf("nil listener - %s:%d %s\n", frame.File, frame.Line, frame.Function)
	}
	rd.Register(l)

	return rd.Data
}

// String implements Zephyr.String() string,
// and is used internally by the HTML renderer
func (str LiveString) string(l Listener) string {
	rd := str()
	rd.Register(l)
	return rd.Data.(string)
}

// LiveInt is the LiveData implementation
// for the `string` type.
type LiveInt func() *DataDep

// NewLiveString returns a "live" string (reactive type LiveString)
func (c *BaseComponent) NewLiveInt(data int) LiveInt {
	// create a new DataDep
	rd := NewDep(data)
	rdPtr := &rd
	// return func type with getter
	rdGetter := LiveInt(func() *DataDep {
		return rdPtr
	})
	return rdGetter
}

// Set implements LiveData.Set(interface{}),
// and is used to set and notify listeners.
func (i LiveInt) Set(newData interface{}) {
	val, ok := newData.(int)
	if !ok {
		panic("invalid data type - fixme")
	}
	// setter func?
	rd := i()
	rd.Data = val
	rd.Notify()
}

func (i LiveInt) Value(l Listener) interface{} {
	rd := i()
	if l == nil {
		// pc := make([]uintptr, 15)
		// n := runtime.Callers(2, pc)
		// frames := runtime.CallersFrames(pc[:n])
		// frame, _ := frames.Next()
		// fmt.Printf("nil listener - %s:%d %s\n", frame.File, frame.Line, frame.Function)
	}
	rd.Register(l)

	return rd.Data
}

// String implements Zephyr.String() string,
// and is used internally by the HTML renderer
func (i LiveInt) string(l Listener) string {
	rd := i()
	rd.Register(l)
	strVal := strconv.Itoa(rd.Data.(int))
	return strVal
}

// LiveBool is the LiveData implementation
// for the `string` type.
type LiveBool func() *DataDep

// NewLiveString returns a "live" string (reactive type LiveString)
func (c *BaseComponent) NewLiveBool(data bool) LiveBool {
	// create a new DataDep
	rd := NewDep(data)
	rdPtr := &rd
	// return func type with getter
	rdGetter := LiveBool(func() *DataDep {
		return rdPtr
	})
	return rdGetter
}

// Set implements LiveData.Set(interface{}),
// and is used to set and notify listeners.
func (b LiveBool) Set(newData interface{}) {
	val, ok := newData.(bool)
	if !ok {
		panic("invalid data type - fixme")
	}
	// setter func?
	rd := b()
	rd.Data = val
	rd.Notify()
}

func (b LiveBool) Value(l Listener) interface{} {
	rd := b()
	if l == nil {
		// pc := make([]uintptr, 15)
		// n := runtime.Callers(2, pc)
		// frames := runtime.CallersFrames(pc[:n])
		// frame, _ := frames.Next()
		// fmt.Printf("nil listener - %s:%d %s\n", frame.File, frame.Line, frame.Function)
	}
	rd.Register(l)

	return rd.Data
}

// String implements Zephyr.String() string,
// and is used internally by the HTML renderer
func (b LiveBool) string(l Listener) string {
	rd := b()
	rd.Register(l)
	bool := rd.Data.(bool)
	if bool {
		return "true"
	}
	return "false"
}

// LiveArr is the LiveData implementation
// for the `string` type.
type LiveArray func() *DataDep

// NewLiveString returns a "live" string (reactive type LiveString)
func (c *BaseComponent) NewLiveArray(data interface{}) LiveArray {
	// create a new DataDep
	rd := NewDep(data)
	rdPtr := &rd

	switch data.(type) {
	case []int, []bool, []string, []float32, []float64, []struct{}, []interface{}, []LiveStructIface:
		// return func type with getter
		rdGetter := LiveArray(func() *DataDep {
			return rdPtr
		})
		return rdGetter
	default:
		panic("error, unsupported array type")
	}
}

// Set implements LiveData.Set(interface{}),
// and is used to set and notify listeners.
func (arr LiveArray) Set(newData interface{}) {
	// val, ok := newData.(int)
	// if !ok {
	// 	panic("invalid data type - fixme")
	// }
	// setter func?
	rd := arr()
	switch newData.(type) {
	case []int, []string, []bool, []uint, []float32, []float64, []rune, []struct{}, []interface{}, []LiveStructIface:
		rd.Data = newData
	}
	rd.Notify()
}

func (arr LiveArray) Value(l Listener) interface{} {
	rd := arr()
	if l == nil {
		// pc := make([]uintptr, 15)
		// n := runtime.Callers(2, pc)
		// frames := runtime.CallersFrames(pc[:n])
		// frame, _ := frames.Next()
		// fmt.Printf("nil listener - %s:%d %s\n", frame.File, frame.Line, frame.Function)
	}
	rd.Register(l)

	return rd.Data
}

func (arr LiveArray) Append(val interface{}) {
	// arr.Set(append(arr.Value(nil).([]zephyr.LiveStructIface), &item2))
	rd := arr()
	// rd.Data = append(rd.Data, val)
	switch rd.Data.(type) {
	case []int:
		val, ok := val.(int)
		if !ok {
			panic("type error")
		}
		rd.Data = append(rd.Data.([]int), val)
	case []string:
		val, ok := val.(string)
		if !ok {
			panic("type error")
		}
		rd.Data = append(rd.Data.([]string), val)
	case []bool:
		val, ok := val.(bool)
		if !ok {
			panic("type error")
		}
		rd.Data = append(rd.Data.([]bool), val)
	case []LiveStructIface:
		val, ok := val.(LiveStruct)
		if !ok {
			panic("type error")
		}
		rd.Data = append(rd.Data.([]LiveStructIface), &val)
		rd.Notify()
	}
}

func arrToString(arr interface{}, l Listener) string {
	str := "[ "
	switch arr.(type) {
	case []int:
		for _, item := range arr.([]int) {
			str += strconv.Itoa(item) + " "
		}
	case []string:
		for _, item := range arr.([]string) {
			str += item + " "
		}
	case []bool:
		for _, item := range arr.([]bool) {
			if item {
				str += "true "
				continue
			}
			str += "false "
		}
	case []LiveStructIface:
		// change me?
		for _, item := range arr.([]LiveStructIface) {
			item.Register(l)
			str += "{" + strings.Split(fmt.Sprintf("%+v ", item), "}")[1] + " }"
		}
	default:
		panic("type not supported")
	}

	str += "]"
	return str
}

// String implements Zephyr.String() string,
// and is used internally by the HTML renderer
func (arr LiveArray) string(l Listener) string {
	rd := arr()
	rd.Register(l)
	switch rd.Data.(type) {
	case []string, []int, []bool, []LiveStructIface:
		return arrToString(rd.Data, l)
	}
	return ""
}

// The struct implementation is slightly different
// than the rest; details below.

type LiveStructIface interface {
	Notify()
	Register(l Listener)
}

type LiveStruct struct {
	Listeners map[string]Listener
}

func (s LiveStruct) Notify() {
	for _, l := range s.Listeners {
		l.Update()
	}
}

func (s *LiveStruct) Register(l Listener) {
	if s.Listeners == nil {
		s.Listeners = map[string]Listener{l.Identifier(): l}
		return
	}
	if l != nil {
		s.Listeners[l.Identifier()] = l
	}
}
