package zephyr

import (
	"fmt"
	"reflect"
	"strings"
)

type Component interface {

	// Public API
	Init()
	Render() VNode

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

	props      map[string]ReactiveData
	data       map[string]ReactiveData
	methods    map[string]interface{}
	components map[string]Component

	parentComponent *BaseComponent

	// Listener is notified of any changes
	// to the variables it is listening to
	Listener ComponentListener

	Node *VNode

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

func (c *BaseComponent) getBase() *BaseComponent {
	return c
}

func (c *BaseComponent) setNode(node *VNode) {
	c.Node = node
}

func RenderWrapper(c Component) VNode {
	// initial render and re-renders, cache unchanged components
	base := c.getBase()
	if base.Node != nil {
		if base.CompareData(base.Node) {
			return *base.Node
		}
	}
	node := c.Render()
	c.getBase().Node = &node
	node.Component = c

	return node
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
		base.Listener = ComponentListener{ID: base.interalID}
	}

	return c

}

func (c *BaseComponent) RegisterComponents(components []Component) {
	c.components = make(map[string]Component)
	for i, childIface := range components {
		child := components[i].getBase()
		onChildUpdate := func() {
			fmt.Println("test")
		}
		child.SetListenerUpdater(onChildUpdate)
		child.getBase().parentComponent = c
		childIface.Init()
		c.components[child.getBase().interalID] = childIface
		// child.getBase().parentComponent = c
	}
}

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

func (c *BaseComponent) GetChildComponent(id string) Component {
	return c.components[id]
}

func (c *BaseComponent) GetMethod(key string) func(args ...interface{}) interface{} {
	if c.methods == nil {
		c.methods = make(map[string]interface{})
	}
	if f, ok := c.methods[key]; ok {
		switch f.(type) {
		case func(), func(...interface{}):
			return func(args ...interface{}) interface{} {
				f.(func(...interface{}))(args)
				return nil
			}
		case func() interface{}, func(...interface{}) interface{}:
			return f.(func(...interface{}) interface{})
		default:
			// fmt.Println(reflect.TypeOf(f))
			return nil
		}
	}
	return nil
}

func (c *BaseComponent) SetMethod(key string, data interface{}) {
	if c.methods == nil {
		c.methods = make(map[string]interface{})
	}
	if _, ok := c.data[key]; ok {
		panic("Method redefined!")
	}

	c.methods[key] = data

	// notify of update
	// c.methods[key].Notify()
}

// func DefineMethods() map[string]MethodFunc {

// }

// // figure out a better name for computed
// func DefineComputed() map[string]ComputedFunc {

// }

// Unused
// func InitComponent(c Component) Component {
// 	// ignoredMethods := []string{"Get", "Render"}
// 	base := c.getBase()
// 	// base := c
// 	// fmt.Println(base)
// 	base.Data = make(map[string]ReactiveData)
// 	structType := reflect.TypeOf(c)

// 	// parse out data from struct fields
// 	for i := 1; i < structType.NumField(); i++ {
// 		field := structType.Field(i)
// 		// fmt.Println(structType.Field(i).Name)
// 		base.Data[field.Name] = newReactiveData(field.Type.Name(), nil)
// 		// fmt.Println(field)
// 	}

// 	// parse out hook and other funcs from struct methods
// 	for i := 0; i < structType.NumMethod(); i++ {
// 		method := structType.Method(i)

// 		switch method.Name {
// 		case "OnInit":
// 			fmt.Println()
// 		case "Get", "Render":
// 		}
// 		// fmt.Println(method.Name)

// 	}

// 	return c
// }
