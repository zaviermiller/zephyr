package main

import (
	"github.com/zaviermiller/zephyr/examples/superbasic/src/components/hello"
	zephyr "github.com/zaviermiller/zephyr/pkg/core"
)

type RootComponent struct {
	zephyr.BaseComponent

	// component access - remove
	HelloComponent zephyr.Component

	counter    int
	setCounter func(int)

	message    string
	setMessage func(string)

	// `vugu:"prop"`

	computedMessage func() string
}

var test = "steamy"

func (rc *RootComponent) Init() {
	// set components to local vars
	rc.HelloComponent = zephyr.NewComponent(&hello.HelloComponent{})

	// var counter = 0
	// rc.ReactiveInt(&counter, "counter")
	// rc.MakeReactive<Int>(&counter, "counter")

	rc.RegisterComponents([]zephyr.Component{
		rc.HelloComponent,
	})

	// rc.counter = 0
	rc.setCounter = rc.BindInt(&rc.counter)

	// rc.message = "Zephyr"
	// rc.setMessage = rc.BindString(&rc.message)

	// rc.computedMessage =

	rc.DefineData(map[string]interface{}{
		"recipient": "Zephyr",
		"messageComputed": func() interface{} {
			return rc.Get("message").(string) + " and world!"
		},
	})

}

func (rc *RootComponent) increaseCounter() {
	rc.setCounter(rc.counter + 1)
	test += "new"
}

// Render is a function that must be implemented by all
// components and is responsible for building the zephyr of the
// component.
func (rc *RootComponent) Render() zephyr.VNode {
	// localVar := reflect.TypeOf(&hello.HelloComponent{}).String()

	return zephyr.BuildElem("div", nil, []zephyr.VNode{
		zephyr.BuildElem("button", map[string]interface{}{
			"onclick":      rc.increaseCounter,
			"ontouchstart": rc.increaseCounter,
		}, []zephyr.VNode{
			zephyr.BuildText(&test),
		}),
		zephyr.BuildComment("OH NO MY MELONS"),
		zephyr.BuildElem("span", nil, []zephyr.VNode{zephyr.BuildText(&test)}),
		rc.HelloComponent.Render(),
		// rc.GetChildComponentWithProps(strings.Split(localVar, ".")[1]).Render(),
		// zephyr.BuildComponent(rc.)
		// zephyr.BuildElem("input", map[string]interface{}{"type": "text", "onchange": func(el js.alue) { rc.Set("message", el.Get("value").String()) }}, nil)
	})
}

// ac.Set("counter", 0)

// func (rc RootComponent) SomeComputedProp() string {
// 	return "Current count: " + strconv.Itoa(rc.Get("counter").(int))
// }
