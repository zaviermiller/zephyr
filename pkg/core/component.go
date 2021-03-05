package zephyr

import (
	"fmt"
	"reflect"
	"runtime"
	"strconv"

	"github.com/zaviermiller/zephyr/pkg/core/vdom"
)

type Component interface {

	// Public API
	Init()
	Render() vdom.VNode

	// Base functions
	Get(string) interface{}
	Set(string, interface{}) interface{}

	// internal use (maybe unnecessary)
	CreateListener(ComponentListener)
	getBase() *BaseComponent
}

// HookFunc is the type used for the hook functions
// that run at certain points in the runtime process.
// HookFuncs should not take or return anything
type HookFunc func()

// MethodFunc may manipulate component state
type MethodFunc func(c *Component, params ...interface{})

// This probably will only allow one return value, is there
// a use case where this doesnt work??
type ComputedFunc func() interface{}

type BaseComponent struct {
	data    map[string]ReactiveData
	methods map[string]MethodFunc

	// ComponentListener is notified of any changes
	// to the variables it is listening to
	Listener *ComponentListener

	// Hooks =-=-=
	// These functions will be called according to
	// the following rules:
	//		Before component is instantiated | BeforeInit() ???
	//		Component is instantiated 			 | OnInit()
	//		Component is mounted to the DOM  | OnMount()
	//		Component is updated 						 | OnUpdate()
	OnInit   HookFunc
	OnMount  HookFunc
	OnUpdate HookFunc
}

func (c *BaseComponent) CreateListener(listener ComponentListener) {
	c.Listener = &listener
}

func (c *BaseComponent) getBase() *BaseComponent {
	return c
}

func (c *BaseComponent) DefineData(dataDefinitions map[string]interface{}) {
	if c.data == nil {
		c.data = make(map[string]ReactiveData)
	}
	for key, val := range dataDefinitions {
		switch val.(type) {
		case func() interface{}:
			c.Set(key, ComputedFunc(val.(func() interface{})))
		default:
			c.Set(key, val)
		}
	}
}

func recurConvert(d interface{}) string {
	switch d.(type) {
	case string:
		return d.(string)
	case int, int8, int16, int32, int64:
		return strconv.Itoa(d.(int))
	case uint, uint8, uint16, uint32, uint64:
		return strconv.Itoa(int(d.(uint)))
	case ComputedFunc:
		return recurConvert(d.(ComputedFunc)())
	}
	return ""

}

// Get is the
func (c *BaseComponent) Get(key string) interface{} {
	pc, _, _, ok := runtime.Caller(1)
	details := runtime.FuncForPC(pc)
	if ok && details != nil {
		fmt.Printf("called from %s\n", details.Name())
	}

	if rd, ok := c.data[key]; ok {
		rd.Register(*c.Listener)
		c.data[key] = rd

		switch rd.Data.(type) {
		// rip no generics
		case func() string:
			return rd.Data.(func() string)()
		case func() int:
			return rd.Data.(func() int)()
		case func() interface{}:
			return rd.Data.(func() interface{})()
		default:
			fmt.Println(reflect.TypeOf(rd.Data))
		}
		return rd.Data
	}
	return nil
}

func (c *BaseComponent) GetStr(key string) string {
	if rd, ok := c.data[key]; ok {
		rd.Register(*c.Listener)
		data := rd.Data
		// immutability smh
		c.data[key] = rd
		return recurConvert(data)
	}
	return ""
}

func (c *BaseComponent) Set(key string, data interface{}) interface{} {
	var newData ReactiveData
	if _, ok := c.data[key]; ok {
		// change to switch
		if reflect.TypeOf(data).Name() != c.data[key].Type {
			panic("Type mismatch!")
		}
		// if _, ok := data.(c.data[key].Type)
		newData = c.data[key]
		newData.Data = data
	} else {
		newData = newReactiveData(reflect.TypeOf(data).Name(), data)
	}

	c.data[key] = newData

	// notify of update
	c.data[key].Notify()

	return newData.Data
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
