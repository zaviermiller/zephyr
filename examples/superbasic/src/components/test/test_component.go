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
	zString zephyr.LiveData
	bigArr  zephyr.LiveData
	// arrLength func() interface{}

	// data tests
	// longArr LiveData

	// Other data
}

func (tc *TestComponent) Init() {
	// tc.BindProp("propArr", tc.longArr)
	tc.BindProp("prop", &tc.zString)

	tc.bigArr = tc.NewLiveArray([]int{1, 2, 3, 4, 5, 6})

	// fmt.Println(tc.longArr)
}

func (tc *TestComponent) calculatedStyle(l *zephyr.VNodeListener) interface{} {
	if len(tc.zString.Value(l).(string)) >= 20 {
		return "color: red;"
	} else {
		return "color: black;"
	}
}

func (tc *TestComponent) calculatedClass(l *zephyr.VNodeListener) interface{} {
	if len(tc.zString.Value(l).(string)) > 30 {
		return "hidden"
	} else {
		return ""
	}
}

func (tc *TestComponent) Render() *zephyr.VNode {
	// fmt.Println(reflect.TypeOf(tc.longArr).String())
	return zephyr.Element("div", nil, []*zephyr.VNode{
		zephyr.Element("p", map[string]interface{}{
			"style": tc.calculatedStyle,
			"class": tc.calculatedClass,
		}, []*zephyr.VNode{
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
