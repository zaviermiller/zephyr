package main

import (
	"reflect"
	"strings"

	"github.com/zaviermiller/zephyr/examples/superbasic/src/components/hello"
	zephyr "github.com/zaviermiller/zephyr/pkg/core"
	"github.com/zaviermiller/zephyr/pkg/core/vdom"
)

type RootComponent struct {
	*zephyr.BaseComponent
}

func (rc *RootComponent) increaseCounter() {
	rc.Set("counter", rc.Get("counter").(int)+1)
}

func (rc *RootComponent) Init() {

	rc.RegisterComponents([]zephyr.Component{
		zephyr.NewComponent(&hello.HelloComponent{&zephyr.BaseComponent{}}),
	})

	rc.DefineData(map[string]interface{}{
		"message": "Hello Zephyr",
		"counter": 0,
		"messageComputed": func() interface{} {
			return rc.Get("message").(string) + " and world!"
		},
	})

}

// Render is a function that must be implemented by all
// components and is responsible for building the vdom of the
// component.
func (rc *RootComponent) Render() vdom.VNode {
	localVar := reflect.TypeOf(&hello.HelloComponent{}).String()
	return vdom.BuildElem("div", nil, []vdom.VNode{
		vdom.BuildElem("button", map[string]interface{}{
			"onclick":      rc.increaseCounter,
			"ontouchstart": rc.increaseCounter,
		}, []vdom.VNode{
			vdom.BuildText("Click me"),
		}),
		vdom.BuildComment("OH NO MY MELONS"),
		vdom.BuildElem("span", nil, []vdom.VNode{vdom.BuildText(rc.GetStr("counter"))}),
		rc.GetChildComponent(strings.Split(localVar, ".")[1]).Render(),
		// rc.GetChildComponentWithProps(strings.Split(localVar, ".")[1]).Render(),
		// vdom.BuildComponent(rc.)
		// vdom.BuildElem("input", map[string]interface{}{"type": "text", "onchange": func(el js.alue) { rc.Set("message", el.Get("value").String()) }}, nil)
	})
}

// ac.Set("counter", 0)

// func (rc RootComponent) SomeComputedProp() string {
// 	return "Current count: " + strconv.Itoa(rc.Get("counter").(int))
// }
