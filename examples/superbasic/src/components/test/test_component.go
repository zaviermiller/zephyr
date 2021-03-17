package test

import (
	"fmt"

	"github.com/zaviermiller/zephyr/examples/superbasic/src/components/text_field"
	zephyr "github.com/zaviermiller/zephyr/pkg/core"
)

var Component = zephyr.NewComponent(&TestComponent{})

type TestComponent struct {
	zephyr.BaseComponent
	ComponentData
}

type ComponentData struct {
	// Props
	longArr   *[]int
	arrLength func() interface{}

	// Other data
}

func (tc *TestComponent) Init() {
	tc.longArr = tc.BindProp("propArr").(*[]int)
	tc.arrLength = tc.BindProp("propComputed").(func() interface{})
	fmt.Println(tc.longArr)
}

func (tc *TestComponent) Render() *zephyr.VNode {
	// fmt.Println(reflect.TypeOf(tc.longArr).String())
	return zephyr.Element("div", nil, []*zephyr.VNode{
		zephyr.Element("p", nil, []*zephyr.VNode{
			zephyr.DynamicText(tc.longArr),
		}),
		tc.ChildComponent(text_field.Component, map[string]interface{}{"initial": tc.arrLength}), /*.BindEvent("change", func(e zephyr.DOMEvent) { tc.setName(e.Target.Value) })*/
	})
}
