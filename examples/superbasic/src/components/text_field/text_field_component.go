package text_field

import (
	zephyr "github.com/zaviermiller/zephyr/pkg/core"
)

var Component = zephyr.NewComponent(&TextFieldComponent{})

type TextFieldComponent struct {
	zephyr.BaseComponent
	Data
}

type Data struct {
	value func() interface{}
}

func (c *TextFieldComponent) Init() {
	c.value = c.BindProp("initial").(func() interface{})
}

func (c *TextFieldComponent) Render() *zephyr.VNode {
	return zephyr.Element("input", map[string]interface{}{"type": "text", "value": c.value}, nil) /*.BindEvent("change", func(e zephyr.DOMEvent) { tc.setName(e.Target.Value) })*/
}

// zephyr.Element("")
