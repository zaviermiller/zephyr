package main

import (
	zephyr "github.com/zaviermiller/zephyr/pkg/core"
	"github.com/zaviermiller/zephyr/pkg/core/vdom"
)

// import core/reactivity

type RootComponent struct {
	// Extend Component struct
	*zephyr.BaseComponent
}

func (rc *RootComponent) Init() {

	// rc.RegisterComponents([]zephyr.Component{
	// 	HomeComponent
	// })

	// define data here (can also be set elsewhere)
	rc.DefineData(map[string]interface{}{
		"message": "Hello Zephyr",
		"counter": 0,
		"computedProp": func() string {
			return rc.Get("message").(string) + " and world!"
		},
	})

	// rc.Methods = zephyr.DefineMethods({
	//   "methodOne": rc.methodOne
	// })
}

// Render is a function that must be implemented by all
// components and is responsible for building the vdom of the
// component.
func (rc *RootComponent) Render() vdom.VNode {
	return vdom.BuildElem("div", nil, []vdom.VNode{
		vdom.BuildElem("button", map[string]interface{}{
			"onclick": func() {
				rc.Set("counter", rc.Get("counter").(int)+1)
			},
		}, []vdom.VNode{
			vdom.BuildText("Click me"),
		}),
		vdom.BuildElem("span", nil, []vdom.VNode{vdom.BuildText(rc.GetStr("counter"))}),
		vdom.BuildElem("h3", nil, []vdom.VNode{vdom.BuildText(rc.Get("computedProp").(string))}),
		// vdom.BuildElem("input", map[string]interface{}{"type": "text", "onchange": func(el js.V	alue) { rc.Set("message", el.Get("value").String()) }}, nil)
	})
}

// ac.Set("counter", 0)

// func (rc RootComponent) SomeComputedProp() string {
// 	return "Current count: " + strconv.Itoa(rc.Get("counter").(int))
// }
