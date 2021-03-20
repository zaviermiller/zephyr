package test

import (
	zephyr "github.com/zaviermiller/zephyr/pkg/core"
)

var Component = zephyr.NewComponent(&TestComponent{})

type TestComponent struct {
	zephyr.BaseComponent
	ComponentData
}

type ComponentData struct {
	// Props
	zString   zephyr.ZephyrData
	arrLength func() interface{}

	// data tests
	// longArr ZephyrData

	// Other data
}

func (tc *TestComponent) Init() {
	// tc.BindProp("propArr", tc.longArr)
	tc.BindProp("prop", &tc.zString)
	// fmt.Println(tc.longArr)
}

func (tc *TestComponent) Render() *zephyr.VNode {
	// fmt.Println(reflect.TypeOf(tc.longArr).String())
	return zephyr.Element("div", nil, []*zephyr.VNode{
		zephyr.Element("p", nil, []*zephyr.VNode{
			zephyr.DynamicText(tc.zString),
		}),
		// tc.ChildComponent(text_field.Component, map[string]interface{}{"initial": tc.arrLength}), /*.BindEvent("change", func(e zephyr.DOMEvent) { tc.setName(e.Target.Value) })*/
	})
}

// tc.Element("div", nil, []*zephyr.VNode {
// 	tc.Element("p", nil, []*zephyr.VNode{
// string!
// 		tc.DynamicText(longArr)
// 	})
// })
