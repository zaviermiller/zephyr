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

// func ()

type VNodeListener struct {
	// Updater func()
	id   string
	node *VNode
	// depTypes []ListenerType
}

func (l VNodeListener) Update() {
	// re-render component on update
	// if l.Updater != nil {
	// 	l.Updater()
	// }
}

func (l VNodeListener) Identifier() string {
	return l.id
}

type VNAttrListener struct {
	// Updater func()
	id   string
	node *VNode
	// depTypes []ListenerType
}

func (l VNAttrListener) Update() {
	// re-render component on update
	// if l.Updater != nil {
	// 	l.Updater()
	// }
	// re-parse attrs, check diffs, render
	updatedAttrs, err := l.node.parseAttrs()
	if err != nil {
		panic(err.Error())
	}
	l.node.RenderChan <- DOMUpdate{Operation: UpdateAttrs, ElementID: l.node.GetDOMSelector(), Data: updatedAttrs}
	// for k, v := range l.node.Attrs {
	// 	newV, ok := updatedAttrs[k]
	// 	if !ok {
	// 		// remove attr
	// 		l.node.RenderChan <- DOMUpdate{Operation: RemoveAttr, ElementID: l.node.DOM_ID, Data: /* some data? */ "test" }
	// 		return
	// 	}
	// 	if newV != v {
	// 		l.node.RenderChan <- DOMUpdate{Operation: UpdateAttr, ElementID: l.node.DOM_ID, Data: newV}
	// 	}
	// set arr
	// for _, newVal := range node.HTMLNode.Attr {
	// 	if newVal.Key == val.Key {
	// 		if newVal.Val == val.Val {
	// 			break
	// 		} else {
	// 			// fmt.Println("mismatched attr, sending update: ", newVal.Val, val.Val)
	// 			// z.QueueUpdate(UpdateAttr, node.DOM_ID, newVal)
	// 			break
	// 		}
	// 	}
	// }
}

func (l VNAttrListener) Identifier() string {
	return l.id
}

type VNContentListener struct {
	// Updater func()
	id   string
	node *VNode
	// depTypes []ListenerType
}

func (l VNContentListener) Update() {
	// re-render component on update
	// if l.Updater != nil {
	// 	l.Updater()
	// }
	updatedContent, err := l.node.parseContent()
	if err != nil {
		panic(err.Error())
	}
	if updatedContent != l.node.Tag {
		l.node.Parent.RenderChan <- DOMUpdate{Operation: UpdateContent, ElementID: l.node.Parent.GetDOMSelector(), Data: updatedContent}
		l.node.Tag = updatedContent
	}
}

func (l VNContentListener) Identifier() string {
	return l.id
}

type VNPropListener struct {
	// Updater func()
	id   string
	node *VNode
	// depTypes []ListenerType
}

func (l VNPropListener) Update() {
	// re-render component on update
	// if l.Updater != nil {
	// 	l.Updater()
	// }
}

func (l VNPropListener) Identifier() string {
	return l.id
}

type VNCalculatorListener struct {
	// Updater func()
	id   string
	node *VNode
	// depTypes []ListenerType
}

func (l VNCalculatorListener) Update() {
	// re-render component on update
	// if l.Updater != nil {
	// 	l.Updater()
	// }
}

func (l VNCalculatorListener) Identifier() string {
	return l.id
}

type VNConditionalListener struct {
	// Updater func()
	id   string
	node *VNode
	// depTypes []ListenerType
}

func (l VNConditionalListener) Update() {
	// re-render component on update
	// if l.Updater != nil {
	// 	l.Updater()
	// }
	l.node.parseConditional()
	// fmt.Println("fc: ", l.node.FirstChild)
	fmt.Println(l.node.GetDOMSelector())
	if l.node.ConditionUpdated {
		if l.node.FirstChild == nil {
			l.node.RenderChan <- DOMUpdate{Operation: Delete, ElementID: l.node.GetDOMSelector(), Data: "self"}
			return
		}
		l.node.RenderChan <- DOMUpdate{Operation: UpdateConditional, ElementID: l.node.GetDOMSelector(), Data: l.node.FirstChild}
		l.node.Tag = l.node.FirstChild.Tag
		var eventRecur func(*VNode)
		eventRecur = func(node *VNode) {
			node.RenderChan = l.node.RenderChan
			if node.events != nil {
				fmt.Println(node.DOM_ID, node.events)
				// l.node.RenderChan <- DOMUpdate{Operation: AddEventListeners, ElementID: node.GetDOMSelector(), Data: node.events}
				l.node.RenderChan <- DOMUpdate{Operation: AddEventListeners, ElementID: node.GetDOMSelector(), Data: node}
			}
			for c := node.FirstChild; c != nil; c = c.NextSibling {
				eventRecur(c)
			}
		}
		eventRecur(l.node)
	}
}

func (l VNConditionalListener) Identifier() string {
	return l.id
}

type VNIteratorListener struct {
	// Updater func()
	id   string
	node *VNode
	// depTypes []ListenerType
}

func (l VNIteratorListener) Update() {
	// re-render component on update
	// if l.Updater != nil {
	// 	l.Updater()
	// }
	if l.node.NodeType != IterativeNode {
		return
	}
	newKeys := l.node.getNewKeys()
	keys := l.node.Keys

	// var bb bytes.Buffer
	// RenderHTML(&bb, l.node)
	// // elId = currentUpdate.Data.(*VNode).DOM_ID
	// renderedHTML := string(bb.Bytes())
	// l.node.RenderChan <- DOMUpdate{Operation: RemoveElements, ElementID: l.node.DOM_ID, Data: renderedHTML}
	switch keyDiff := len(newKeys) - len(keys); {
	case keyDiff == 0:
		// check for reorder
		for i, val := range newKeys {
			if oldI := indexOf(keys, val); oldI == -1 {
				// deleted and replaced
			} else if oldI != i {
				l.node.RenderChan <- DOMUpdate{Operation: SwapChildren, ElementID: l.node.Parent.GetDOMSelector(), Data: [2]int{i, oldI}}
				swap(keys, i, oldI)
			}
		}
	case keyDiff > 0:
		// item was added
		curr := l.node.FirstChild
		for i, val := range newKeys {
			if oldI := indexOf(keys, val); oldI == -1 {
				// new element
				if keyDiff <= 0 {
					// element deleted and replaced
					// todo
				} else {
					// element inserted at i
					keyDiff--
					keys = insertAt(keys, i, val)
					newNode := l.node.IterRender(i, l.node.Content.(LiveArray)().Data.([]LiveStruct)[i])
					// newNode.DOM_ID = l.node.DOM_ID + "-k[" + newNode.key.(string) + "]"
					newNode.DOM_ID = l.node.DOM_ID
					newNode.Attrs["data-key"] = newNode.key
					if curr == nil {
						// appended
						prev := l.node.LastChild
						prev.NextSibling = newNode
						newNode.PrevSibling = prev
						l.node.LastChild = newNode
						i := 0
						cTest := newNode.PrevSibling
						for cTest != nil {
							i++
							cTest = cTest.PrevSibling
						}

						l.node.RenderChan <- DOMUpdate{Operation: InsertAfter, ElementID: prev.GetDOMSelector(), Data: newNode}
					} else {
						prev := curr.PrevSibling
						if prev != nil {
							prev.NextSibling = newNode
						}
						curr.PrevSibling = newNode
						newNode.PrevSibling = prev
						newNode.NextSibling = curr
						i := 0
						cTest := newNode.PrevSibling
						for cTest != nil {
							i++
							cTest = cTest.PrevSibling
						}

						l.node.RenderChan <- DOMUpdate{Operation: InsertBefore, ElementID: curr.GetDOMSelector(), Data: newNode}
					}
					// if curr.NextSibling == nil {
					// 	l.node.RenderChan <- DOMUpdate{Operation: InsertAfter, ElementID: curr.PrevSibling.DOM_ID, Data: curr}
					// 	continue
					// }
				}
			} else if oldI != i {
				// swap the elements
				l.node.RenderChan <- DOMUpdate{Operation: SwapChildren, ElementID: l.node.Parent.GetDOMSelector(), Data: [2]int{i, oldI}}
				swap(keys, i, oldI)
			}
			if curr != nil {
				curr = curr.NextSibling
			}
		}
	case keyDiff < 0:
		// item was removed
		curr := l.node.FirstChild

		for i, val := range keys {
			if oldI := indexOf(newKeys, val); oldI == -1 {
				if keyDiff >= 0 {
					// replaced
				} else {
					// remove the attr
					l.node.RenderChan <- DOMUpdate{Operation: Delete, ElementID: curr.GetDOMSelector(), Data: "before"}
					keys = append(keys[:i], keys[:i+1])
					keyDiff++
				}
			} else if oldI != i {
				//swap?
				l.node.RenderChan <- DOMUpdate{Operation: SwapChildren, ElementID: l.node.Parent.GetDOMSelector(), Data: [2]int{i, oldI}}
				swap(keys, i, oldI)
			}
		}
	}
	l.node.Keys = newKeys
	// fmt.Println("received iterator update")
	// l.node.parseConditional()
	// // fmt.Println("fc: ", l.node.FirstChild)
	// if l.node.ConditionUpdated {
	// 	l.node.RenderChan <- DOMUpdate{Operation: UpdateConditional, ElementID: l.node.DOM_ID, Data: l.node.FirstChild}
	// 	var eventRecur func(*VNode)
	// 	eventRecur = func(node *VNode) {
	// 		node.RenderChan = l.node.RenderChan
	// 		if node.events != nil {
	// 			l.node.RenderChan <- DOMUpdate{Operation: AddEventListeners, ElementID: node.DOM_ID, Data: node.events}
	// 		}
	// 		for c := node.FirstChild; c != nil; c = c.NextSibling {
	// 			eventRecur(c)
	// 		}
	// 	}
	// 	eventRecur(l.node)
	// }
}

// arr helper funcs for key comps

func indexOf(arr []interface{}, val interface{}) int {
	for i, item := range arr {
		if item == val {
			return i
		}
	}
	return -1
}

func insertAt(arr []interface{}, index int, val interface{}) []interface{} {
	tmp := append(arr[:index], val)
	return append(tmp, arr[index:])
}

func swap(arr []interface{}, index1, index2 int) []interface{} {
	tmp := arr[index1]
	arr[index1] = arr[index2]
	arr[index2] = tmp
	return arr
}

func (l VNIteratorListener) Identifier() string {
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
func NewLiveString(data string) LiveString {
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
func NewLiveInt(data int) LiveInt {
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
func NewLiveBool(data bool) LiveBool {
	// create a new DataDep
	rd := NewDep(data)
	rdPtr := &rd
	// return func type with getter
	rdGetter := LiveBool(func() *DataDep {
		return rdPtr
	})
	// c.data[rdGetter] = rd
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
	// fmt.Println(rd.Listeners)
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
func NewLiveArray(data interface{}) LiveArray {
	// create a new DataDep
	rd := NewDep(data)
	rdPtr := &rd

	switch data.(type) {
	case []int, []bool, []string, []float32, []float64, []struct{}, []interface{}, []LiveStruct:
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
	case []int, []string, []bool, []uint, []float32, []float64, []rune, []struct{}, []interface{}, []LiveStruct:
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

func (arr LiveArray) At(l Listener, i int) interface{} {
	arrV := arr.Value(l)
	switch arrV.(type) {
	case []LiveStruct:
		return arrV.([]LiveStruct)[i]
	case []string:
		return arrV.([]string)[i]
	case []bool:
		return arrV.([]bool)[i]
	case []int:
		return arrV.([]int)[i]
	default:
		panic("type not supported")
	}
}

func (arr LiveArray) Append(val interface{}) {
	// arr.Set(append(arr.Value(nil).([]zephyr.LiveStruct), &item2))
	rd := arr()
	// rd.Data = append(rd.Data, val)
	switch rd.Data.(type) {
	case []int:
		val, ok := val.(int)
		if !ok {
			panic("[]int type error")
		}
		rd.Data = append(rd.Data.([]int), val)
	case []string:
		val, ok := val.(string)
		if !ok {
			panic("[]string type error")
		}
		rd.Data = append(rd.Data.([]string), val)
	case []bool:
		val, ok := val.(bool)
		if !ok {
			panic("[]bool type error")
		}
		rd.Data = append(rd.Data.([]bool), val)
	case []LiveStruct:
		v, ok := val.(LiveStruct)
		if !ok {
			panic("[]LiveStruct type error")
		}
		rd.Data = append(rd.Data.([]LiveStruct), v)
		rd.Notify()
		fmt.Println(rd.Listeners)
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
	case []LiveStruct:
		// change me?
		for _, item := range arr.([]LiveStruct) {
			item.Register(l)
			structSplit := strings.Split(fmt.Sprintf("%+v ", item), "}")
			printed := structSplit[len(structSplit)-2]
			str += "{" + printed + " } "
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
	case []string, []int, []bool, []LiveStruct:
		return arrToString(rd.Data, l)
	}
	return ""
}

// The struct implementation is slightly different
// than the rest; details below.

type LiveStruct interface {
	Notify()
	Register(l Listener)
}

type LiveStructImpl struct {
	Listeners map[string]Listener
}

func (s LiveStructImpl) Notify() {
	for _, l := range s.Listeners {
		l.Update()
	}
}

func (s *LiveStructImpl) Register(l Listener) {
	if s.Listeners == nil {
		s.Listeners = map[string]Listener{l.Identifier(): l}
		return
	}
	if l != nil {
		s.Listeners[l.Identifier()] = l
	}
}
