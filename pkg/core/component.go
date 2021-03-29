package zephyr

import (
	"strconv"
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
	props map[string]interface{}
	// reactive data - internal use. check reactivity.go
	// data map[LiveData]DataDep
	// methods map[string]interface{}
	// components map[string]Component

	parentComponent *BaseComponent

	// Listener is notified of any changes
	// to the variables it is listening to
	// Listener ComponentListener

	// Node is a reference to the components
	// root vNode
	Node *VNode

	// Context is a reference to the ZephyrApp
	// which is necessary to register global events
	// and stuff
	// Context *ZephyrApp

	// Hooks =-=-=
	// These functions will be called according to
	// the following rules:
	//		Before component is instantiated | BeforeInit() ???
	//		Component is instantiated 			 | OnInit()
	//		Component is mounted to the DOM  | OnMount()
	//		Component is updated 						 | OnUpdate()
}

// ChildComponent calls the render func of a child component
func (parent *BaseComponent) Component(c Component) *VNode {
	// set context based on parent
	updateQ := parent.Node.RenderChan
	base := c.getBase()
	InitWrapper(c)

	// render component
	RenderWrapper(c, updateQ)
	// base.Node.Content = c
	base.Node.Component = true
	base.parentComponent = parent
	return base.Node
}

// ChildComponent calls the render func of a child component
func (parent *BaseComponent) ComponentWithProps(c Component, props map[string]interface{}) *VNode {
	// set context based on parent
	updateQ := parent.Node.RenderChan
	base := c.getBase()

	// parse and pass props
	base.props = props
	InitWrapper(c)

	// render component
	RenderWrapper(c, updateQ)
	// base.Node.Content = c
	base.Node.Component = true
	base.parentComponent = parent
	return base.Node
}

func (c *BaseComponent) BindProp(propName string) interface{} {
	base := c.getBase()
	val, ok := base.props[propName]
	if !ok {
		panic("Zephyr framework error")
	}

	// type safety!!
	switch val.(type) {
	case LiveArray:
		return val
	case LiveBool:
		return val
	case LiveString:
		return val
	case LiveInt:
		return val
	case LiveStruct:
		return val
	// TODO
	case func(Listener) interface{}:
		evald := val.(func(Listener) interface{})(base.Node.GetOrCreateListener("prop"))
		return evald
		// add evald to props for comp later
		base.props[propName+"Evaluated"] = evald
	default:
		panic("props must be live data or a calculator")
	}
	return nil
}

// TODO
func propCalcUpdater(node *VNode) {
	if !node.Component {
		panic("props can only be updated on components")
	}
	base := node.Content.(Component).getBase()
	for key, val := range base.props {
		switch val.(type) {
		case func(Listener) interface{}:
			evald := val.(func(Listener) interface{})(base.Node.GetOrCreateListener("prop"))
			if evald != base.props[key+"Evaluated"] {
				// node.RenderChan <- DOMUpdate{Operation: }

			}
		}
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
	// base := c.getBase()
	c.Init()
}

func RenderWrapper(c Component, updateQueue chan DOMUpdate) *VNode {
	// initial render and re-renders, cache unchanged components
	base := c.getBase()
	// need da q
	base.Node = &VNode{RenderChan: updateQueue}
	node := c.Render()
	base.Node = node
	i := 0
	cTest := base.Node.PrevSibling
	for cTest != nil {
		i++
		cTest = cTest.PrevSibling
	}

	base.Node.DOM_ID = "zComponent" + strconv.Itoa(i)

	// fixme
	var recurChanSet func(node *VNode)
	recurChanSet = func(node *VNode) {
		if node == nil {
			return
		}
		node.RenderChan = updateQueue
		for c := node.FirstChild; c != nil; c = c.NextSibling {
			// if c.NodeType == ConditionalNode {
			// 	c.parseConditional()
			// }
			recurChanSet(c)
		}
	}
	recurChanSet(base.Node)
	return node
}

// func recurSetListenerUpdater(node *VNode, ctx *ZephyrApp) {
// 	if node == nil || node.Listener == nil {
// 		return
// 	}
// 	node.Listener.Updater = func() {
// 		if ctx.RootNode != nil {
// 			go ctx.CompareNode(node)
// 			return
// 		}
// 		fmt.Println("error", ctx.RootNode)
// 	}
// 	curr := node.FirstChild
// 	for curr != nil {
// 		recurSetListenerUpdater(curr, ctx)
// 		curr = curr.NextSibling
// 	}
// }

func UpdateWrapper(c Component) {
	//
}

// mayhaps
// func NewComponentWithProps()

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
// 	c.props = make(map[string]DataDep)
// 	for key, val := range propDefs {
// 		c.Set(key, val)
// 	}
// }

// DefineData is a wrapper that initializes and creates the components
// data map from an input
// func (c *BaseComponent) DefineData(dataDefs map[string]interface{}) {
// 	c.data = make(map[string]DataDep)
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
// 		c.data = make(map[string]DataDep)
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
// 		c.data = make(map[string]DataDep)
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
// 		c.data = make(map[string]DataDep)
// 	}
// 	var newData DataDep
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
// 			newData = newDataDep("Computed", ComputedFunc(data.(func() interface{})))
// 		default:
// 			newData = newDataDep(reflect.TypeOf(data).String(), data)
// 		}
// 	}

// 	// fmt.Println("set: ", key)

// 	c.data[key] = newData

// 	// notify of update
// 	c.data[key].Notify()

// 	return newData.Data
// }
