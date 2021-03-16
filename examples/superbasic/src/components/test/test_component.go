package test

import (
	"github.com/zaviermiller/zephyr/examples/superbasic/src/components/text_field"
	zephyr "github.com/zaviermiller/zephyr/pkg/core"
)

var Component = zephyr.NewComponent(&TestComponent{})

type TestComponent struct {
	zephyr.BaseComponent
	ComponentData
}

type ComponentData struct {
	name    string
	setName func(string)
}

func (tc *TestComponent) Init() {
	tc.name = "test"
	tc.setName = tc.BindString(&tc.name)
}

func (tc *TestComponent) Render() zephyr.VNode {
	return *zephyr.Element("div", nil, []*zephyr.VNode{
		zephyr.Element("p", nil, []*zephyr.VNode{
			zephyr.DynamicText(&tc.name),
		}),
		zephyr.ChildComponent(text_field.Component), /*.BindEvent("change", func(e zephyr.DOMEvent) { tc.setName(e.Target.Value) })*/
	})
}
