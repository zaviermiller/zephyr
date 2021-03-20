package zephyr

import (
	"fmt"
	"reflect"
	"strings"
)

type Component interface {

	// Public API
	Init()
	Render() *VNode

	// Base functions
	// Get(string) interface{}
	// Set(string, interface{}) interface{}
	// GetChildComponent(id string) Component

	// interally used to get the base struct and
	// ensure that user defined components
	getBase() *BaseComponent
}

type Mounter interface {
	OnMount()
}

type Updater interface {
	OnUpdate()
}

// This probably will only allow one return value, is there
// a use case where this doesnt work??
type ComputedFunc func() interface{}

type BaseComponent struct {
	interalID string

	props map[string]interface{}
	// reactive data - internal use. check reactivity.go
	data    map[ZephyrData]ReactiveData
	methods map[string]interface{}
	// components map[string]Component

	parentComponent *BaseComponent

	// Listener is notified of any changes
	// to the variables it is listening to
	Listener ComponentListener

	Node *VNode

	// Context is a reference to the ZephyrApp
	// which is necessary to register global events
	// and stuff
	Context *ZephyrApp

	// Hooks =-=-=
	// These functions will be called according to
	// the following rules:
	//		Before component is instantiated | BeforeInit() ???
	//		Component is instantiated 			 | OnInit()
	//		Component is mounted to the DOM  | OnMount()
	//		Component is updated 						 | OnUpdate()
}

// maybe refactor -- was kinda hacking
func (c *BaseComponent) SetListenerUpdater(f func()) {
	c.Listener.Updater = f
}

func (c *BaseComponent) BindProp(propName string, destLocation interface{}) {
	base := c.getBase()
	val, ok := base.props[propName]
	if !ok {
		panic("Zephyr framework error")
	}

	switch val.(type) {
	case ZephyrData:
		if _, ok := destLocation.(*ZephyrData); !ok {
			panic("prop must be same type as local var")
		}
		ptr := destLocation.(*ZephyrData)
		*ptr = val.(ZephyrData)
		// TODO
		// case func() interface{}:
		// if _, ok := destLocation.(*func() interface); !ok {
		// 	panic("prop must be same type as local var")
		// }
		// ptr := destLocation.(*ZephyrData)
		// *ptr = val.(ZephyrData)

	}

}

func (c *BaseComponent) getBase() *BaseComponent {
	return c
}

func (c *BaseComponent) setNode(node *VNode) {
	c.Node = node
}

// The following functions are wrappers around the hooks,
// which get called at different lifecycle events. The
// wrappers exist to run some code before or after the
// user run code which may be necessary.

func InitWrapper(c Component) {
	base := c.getBase()
	if base.Context == nil {
		base.Context = base.parentComponent.Context
	}
	c.Init()
}

func RenderWrapper(c Component) *VNode {
	// initial render and re-renders, cache unchanged components
	base := c.getBase()
	node := c.Render()
	base.Node = node
	// create listener for root changes
	// listener for attr changes
	// listener for text changes
	// listener for prop changes
	// listener for computed changes
	ListenerFunc := func() {
		// vdom := app.RootComponent.Render()
		// fmt.Println("update detected! new vdom: " + func() string { b, _ := json.Marshal(app.RootNode); return string(b) }())
		go base.Context.CompareDOM(base.Node)
	}
	// move?
	base.SetListenerUpdater(ListenerFunc)

	recurSetListenerUpdater(base.Node, base.Context)

	return node
}

func recurSetListenerUpdater(node *VNode, ctx *ZephyrApp) {
	if node == nil || node.Listener == nil {
		return
	}
	fmt.Println(node)
	node.Listener.Updater = func() {
		fmt.Println("updater called for ", node.DOM_ID)
		go ctx.CompareDOM(node)
	}
	curr := node.FirstChild
	for curr != nil {
		recurSetListenerUpdater(curr, ctx)
		curr = curr.NextSibling
	}
}

func UpdateWrapper(c Component) {
	//
}

func NewComponent(c Component) Component {
	base := c.getBase()

	// create the id so it can be found again
	componentId := strings.Split(reflect.TypeOf(c).String(), ".")
	if len(componentId) != 2 {
		// fmt.Println(componentId)
	} else {
		base.interalID = componentId[1]
		base.Listener = ComponentListener{id: base.interalID}
	}

	return c

}

// func (c *BaseComponent) RegisterComponents(components []Component) {
// 	c.components = make(map[string]Component)
// 	for i, childIface := range components {
// 		child := components[i].getBase()
// 		onChildUpdate := func() {
// 			fmt.Println("test")
// 		}
// 		child.SetListenerUpdater(onChildUpdate)
// 		child.getBase().parentComponent = c
// 		childIface.Init()
// 		c.components[child.getBase().interalID] = childIface
// 		// child.getBase().parentComponent = c
// 	}
// }

// DefineData is a wrapper that initializes and creates the components
// data map from an input
// func (c *BaseComponent) DefineProps(propDefs map[string]interface{}) {
// 	c.props = make(map[string]ReactiveData)
// 	for key, val := range propDefs {
// 		c.Set(key, val)
// 	}
// }

// DefineData is a wrapper that initializes and creates the components
// data map from an input
// func (c *BaseComponent) DefineData(dataDefs map[string]interface{}) {
// 	c.data = make(map[string]ReactiveData)
// 	for key, val := range dataDefs {
// 		c.Set(key, val)
// 	}
// }

// DefineMethods initializes and creates the inputed methods
// func (c *BaseComponent) DefineMethods(methodDefs map[string]interface{}) {
// 	c.methods = map[string]interface{}{}
// 	for key, val := range methodDefs {
// 		c.SetMethod(key, val)
// 	}
// }

// func recurToString(d interface{}) string {
// 	switch d.(type) {
// 	case string:
// 		return d.(string)
// 	case int, int8, int16, int32, int64:
// 		return strconv.Itoa(d.(int))
// 	case uint, uint8, uint16, uint32, uint64:
// 		return strconv.Itoa(int(d.(uint)))
// 	case ComputedFunc:
// 		return recurToString(d.(ComputedFunc)())
// 	}
// 	return ""

// }

// Get is the public function used to get values from the
// components data
// func (c *BaseComponent) Get(key string) interface{} {
// 	if c.data == nil {
// 		c.data = make(map[string]ReactiveData)
// 	}
// 	if rd, ok := c.data[key]; ok {
// 		rd.Register(c.Listener)
// 		c.data[key] = rd

// 		switch rd.Data.(type) {
// 		// rip no generics
// 		case ComputedFunc:
// 			return rd.Data.(ComputedFunc)()
// 		default:
// 			return rd.Data
// 		}
// 	}
// 	return nil
// }

// func (c *BaseComponent) GetStr(key string) string {
// 	if c.data == nil {
// 		c.data = make(map[string]ReactiveData)
// 	}
// 	if rd, ok := c.data[key]; ok {
// 		rd.Register(c.Listener)
// 		data := rd.Data
// 		// immutability smh
// 		c.data[key] = rd
// 		return recurToString(data)
// 	}
// 	return ""
// }

// func (c *BaseComponent) Set(key string, data interface{}) interface{} {
// 	if c.data == nil {
// 		c.data = make(map[string]ReactiveData)
// 	}
// 	var newData ReactiveData
// 	if _, ok := c.data[key]; ok {
// 		// change to switch
// 		if reflect.TypeOf(data).String() != c.data[key].Type || reflect.TypeOf(data).Kind() == reflect.Func {
// 			panic("Type mismatch or computed redefinition")
// 		}
// 		// if _, ok := data.(c.data[key].Type)
// 		newData = c.data[key]
// 		newData.Data = data
// 	} else {
// 		switch data.(type) {
// 		case func() interface{}, func():
// 			newData = newReactiveData("Computed", ComputedFunc(data.(func() interface{})))
// 		default:
// 			newData = newReactiveData(reflect.TypeOf(data).String(), data)
// 		}
// 	}

// 	// fmt.Println("set: ", key)

// 	c.data[key] = newData

// 	// notify of update
// 	c.data[key].Notify()

// 	return newData.Data
// }
